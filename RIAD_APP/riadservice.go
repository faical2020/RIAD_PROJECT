package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"RIAD_APP/internal/db"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func structToMap(obj interface{}) map[string]interface{} {
	data, _ := json.Marshal(obj)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return m
}

func roomsToMaps(rooms []db.Room) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rooms))
	for i, r := range rooms {
		result[i] = structToMap(r)
	}
	return result
}

func reservationsToMaps(reservations []db.Reservation) []map[string]interface{} {
	result := make([]map[string]interface{}, len(reservations))
	for i, r := range reservations {
		result[i] = structToMap(r)
	}
	return result
}

var slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

type RiadService struct {
	ctx               context.Context
	token             string
	apiBaseURL        string
	app               *application.App
	grpcConn          *grpc.ClientConn
	grpcClient        pb.SyncServiceClient
	lastSyncTimestamp int64
}

func NewRiadService() *RiadService {
	return &RiadService{apiBaseURL: "http://localhost:8081/api/v1"}
}

func (s *RiadService) SetApp(app *application.App) { s.app = app }

func (s *RiadService) dialGRPC() {
	if s.grpcConn != nil {
		return
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slogger.Warn("gRPC dial failed", "error", err)
		return
	}
	s.grpcConn = conn
	s.grpcClient = pb.NewSyncServiceClient(conn)
	slogger.Info("gRPC client created")
}

func (s *RiadService) grpcCtx() context.Context {
	md := metadata.Pairs("authorization", "Bearer "+s.token)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func (s *RiadService) SetToken(token string) {
	s.token = token
	slogger.Info("token set", "length", len(token))
	s.dialGRPC()
	go s.performSync()
	go s.startGRPCStream()
}

func (s *RiadService) startGRPCStream() {
	if s.token == "" || s.grpcClient == nil {
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			slogger.Debug("connecting to gRPC sync stream")
			stream, err := s.grpcClient.StreamUpdates(s.ctx, &pb.SyncRequest{Token: s.token})
			if err != nil {
				slogger.Warn("gRPC stream failed, retrying in 5s", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}
			slogger.Info("connected to gRPC sync stream")

			for {
				event, err := stream.Recv()
				if err != nil {
					slogger.Warn("gRPC stream recv error, reconnecting", "error", err)
					break
				}
				s.handleSyncEvent(event)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *RiadService) handleSyncEvent(event *pb.SyncEvent) {
	switch event.Type {
	case pb.SyncEvent_ROOM_UPDATED:
		if room := event.GetRoom(); room != nil {
			slogger.Info("room update from gRPC", "id", room.Id)
			db.SaveRoom(room.Id, int(room.Number), room.Type, room.Price, room.Description, room.Equipments, room.Status, room.CleaningStatus)
			if s.app != nil {
				s.app.Event.Emit("sync:rooms", "updated")
			}
		}
	case pb.SyncEvent_RESERVATION_UPDATED:
		if res := event.GetReservation(); res != nil {
			slogger.Info("reservation update from gRPC", "id", res.Id)
			db.SaveReservation(res.Id, res.UserId, res.RoomId, res.StartDate, res.EndDate, res.Amount, res.Status)
			db.MarkSynced("reservations", res.Id)
			if s.app != nil {
				s.app.Event.Emit("sync:reservations", "updated")
			}
		}
	case pb.SyncEvent_CONSOMMATION_UPDATED:
		if s.app != nil {
			s.app.Event.Emit("sync:consommations", "updated")
		}
	}
}

func (s *RiadService) SetContext(ctx context.Context) { s.ctx = ctx }

func (s *RiadService) GetLocalRooms() ([]map[string]interface{}, error) {
	rooms, err := db.GetRooms()
	if err != nil {
		return nil, err
	}
	return roomsToMaps(rooms), nil
}

func (s *RiadService) GetLocalReservations() ([]map[string]interface{}, error) {
	reservations, err := db.GetReservations()
	if err != nil {
		return nil, err
	}
	return reservationsToMaps(reservations), nil
}

func (s *RiadService) CreateLocalReservation(userID, roomID, start, end string, amount float64) (string, error) {
	if s.grpcClient != nil {
		resp, err := s.grpcClient.CreateReservation(s.grpcCtx(), &pb.CreateReservationRequest{
			UserId: userID, RoomId: roomID, StartDate: start, EndDate: end, Amount: amount,
		})
		if err == nil && resp.Id != "" {
			slogger.Info("reservation created via gRPC", "id", resp.Id)
			db.SaveReservation(resp.Id, userID, roomID, start, end, amount, resp.Status)
			db.MarkSynced("reservations", resp.Id)
			return resp.Id, nil
		}
		slogger.Warn("gRPC CreateReservation failed, saving locally", "error", err)
	}

	res := db.Reservation{
		ID:        uuid.New().String(),
		UserID:    userID,
		ChambreID: roomID,
		DateDebut: start,
		DateFin:   end,
		Montant:   amount,
		Statut:    "pending",
	}
	if err := db.DB.Create(&res).Error; err != nil {
		return "", fmt.Errorf("failed to save local reservation: %w", err)
	}
	return res.ID, nil
}

func (s *RiadService) UpdateLocalRoom(id string, num int, roomType string, price float64, desc, equip, status string) error {
	slogger.Info("updating local room", "id", id, "num", num)
	return db.SaveRoom(id, num, roomType, price, desc, equip, status, "propre")
}

func (s *RiadService) UpdateLocalReservation(id, userId, roomId, start, end string, amount float64, status string) error {
	slogger.Info("updating local reservation", "id", id)
	return db.SaveReservation(id, userId, roomId, start, end, amount, status)
}

func (s *RiadService) UpdateCleaningStatus(id, status string) error {
	slogger.Info("updating cleaning status", "id", id, "status", status)

	if s.grpcClient != nil {
		_, err := s.grpcClient.UpdateCleaningStatus(s.grpcCtx(), &pb.UpdateCleaningStatusRequest{RoomId: id, CleaningStatus: status})
		if err == nil {
			if s.app != nil {
				s.app.Event.Emit("sync:rooms", "updated")
			}
			return nil
		}
		slogger.Warn("gRPC UpdateCleaningStatus failed", "error", err)
	}

	room, err := db.GetRoomByID(id)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if err := db.SaveRoom(id, room.Numero, room.Type, room.Prix, room.Description, room.Equipements, room.Statut, status); err != nil {
		return err
	}

	if s.app != nil {
		s.app.Event.Emit("sync:rooms", "updated")
	}
	return nil
}

func (s *RiadService) MarkAsSynced(table, id string) error {
	return db.MarkSynced(table, id)
}

func (s *RiadService) GetUnsynced(table string) ([]map[string]interface{}, error) {
	reservations, err := db.GetUnsynced(table)
	if err != nil {
		return nil, err
	}
	return reservationsToMaps(reservations), nil
}

func (s *RiadService) StartSyncLoop() {
	go func() {
		slogger.Info("background sync loop started")
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		s.performSync()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.performSync()
			}
		}
	}()
}

func (s *RiadService) performSync() {
	if s.token == "" {
		slogger.Debug("sync skipped: no token")
		return
	}
	slogger.Debug("running gRPC background sync")
	s.gRPCPullUpdates()
	s.gRPCSyncReservations()
}

func (s *RiadService) gRPCPullUpdates() {
	if s.grpcClient == nil {
		slogger.Debug("gRPC pull skipped: no client")
		return
	}

	resp, err := s.grpcClient.SyncData(s.grpcCtx(), &pb.SyncDataRequest{LastSequenceId: s.lastSyncTimestamp})
	if err != nil {
		slogger.Warn("gRPC SyncData failed", "error", err)
		return
	}

	for _, room := range resp.Rooms {
		if err := db.SaveRoom(room.Id, int(room.Number), room.Type, room.Price, room.Description, room.Equipments, room.Status, room.CleaningStatus); err != nil {
			slogger.Error("failed to save pulled room", "id", room.Id, "error", err)
		}
	}

	for _, res := range resp.Reservations {
		if err := db.SaveReservation(res.Id, res.UserId, res.RoomId, res.StartDate, res.EndDate, res.Amount, res.Status); err != nil {
			slogger.Error("failed to save pulled reservation", "id", res.Id, "error", err)
		}
		db.MarkSynced("reservations", res.Id)
	}

	s.lastSyncTimestamp = time.Now().Unix()
	slogger.Info("sync pull complete", "rooms", len(resp.Rooms), "reservations", len(resp.Reservations))
}

func (s *RiadService) gRPCSyncReservations() {
	if s.grpcClient == nil {
		return
	}

	unsynced, err := db.GetUnsynced("reservations")
	if err != nil || len(unsynced) == 0 {
		return
	}

	for _, res := range unsynced {
		resp, err := s.grpcClient.CreateReservation(s.grpcCtx(), &pb.CreateReservationRequest{
			UserId: res.UserID, RoomId: res.ChambreID,
			StartDate: res.DateDebut, EndDate: res.DateFin, Amount: res.Montant,
		})
		if err != nil {
			slogger.Warn("gRPC sync failed for reservation", "id", res.ID, "error", err)
			continue
		}
		if resp.Id != "" {
			db.MarkSynced("reservations", res.ID)
			slogger.Info("synced reservation", "local_id", res.ID, "server_id", resp.Id)
		}
	}
}
