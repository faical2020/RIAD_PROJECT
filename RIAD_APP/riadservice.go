package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"RIAD_APP/internal/db"
	"RIAD_APP/pkg/logic"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"github.com/wailsapp/wails/v3/pkg/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var debugLog *log.Logger

func init() {
	f, err := os.OpenFile("/tmp/riad_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not open debug log: %v\n", err)
		return
	}
	debugLog = log.New(f, "[RIAD_DEBUG] ", log.LstdFlags)
}

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
	return &RiadService{
		apiBaseURL: "http://localhost:8081/api/v1",
	}
}

func (s *RiadService) SetApp(app *application.App) {
	s.app = app
}

func (s *RiadService) dialGRPC() {
	if s.grpcConn != nil {
		return
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("gRPC Dial failed: %v\n", err)
		}
		return
	}
	s.grpcConn = conn
	s.grpcClient = pb.NewSyncServiceClient(conn)
	if debugLog != nil {
		debugLog.Println("gRPC client created successfully")
	}
}

func (s *RiadService) grpcCtx() context.Context {
	md := metadata.Pairs("authorization", "Bearer "+s.token)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func (s *RiadService) SetToken(token string) {
	s.token = token
	if debugLog != nil {
		debugLog.Printf("SetToken called: token received (length: %d)\n", len(token))
	}
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
			if debugLog != nil {
				debugLog.Println("Attempting to connect to gRPC sync stream...")
			}

			stream, err := s.grpcClient.StreamUpdates(s.ctx, &pb.SyncRequest{
				Token: s.token,
			})
			if err != nil {
				if debugLog != nil {
					debugLog.Printf("gRPC StreamUpdates failed: %v. Retrying in 5s...\n", err)
				}
				time.Sleep(5 * time.Second)
				continue
			}

			if debugLog != nil {
				debugLog.Println("Connected to gRPC sync stream!")
			}

			for {
				event, err := stream.Recv()
				if err != nil {
					if debugLog != nil {
						debugLog.Printf("gRPC stream recv error: %v. Reconnecting...\n", err)
					}
					break
				}
				s.handleSyncEvent(event)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *RiadService) handleSyncEvent(event *pb.SyncEvent) {
	if debugLog != nil {
		debugLog.Printf("!!! gRPC EVENT RECEIVED: type=%v, id=%s\n", event.Type, event.EntityId)
	}

	switch event.Type {
	case pb.SyncEvent_ROOM_UPDATED:
		if room := event.GetRoom(); room != nil {
			log.Printf("Updating local room %s from gRPC", room.Id)
			db.SaveRoom(room.Id, int(room.Number), room.Type, room.Price, room.Description, room.Equipments, room.Status, room.CleaningStatus)
			if s.app != nil {
				log.Printf("Emitting Wails event sync:rooms for room %s", room.Id)
				s.app.Event.Emit("sync:rooms", "updated")
			} else {
				log.Printf("ERROR: s.app is nil, cannot emit event")
			}
		}
	case pb.SyncEvent_RESERVATION_UPDATED:
		if res := event.GetReservation(); res != nil {
			log.Printf("Updating local reservation %s from gRPC", res.Id)
			db.SaveReservation(res.Id, res.UserId, res.RoomId, res.StartDate, res.EndDate, res.Amount, res.Status)
			db.MarkSynced("reservations", res.Id)
			if s.app != nil {
				log.Printf("Emitting Wails event sync:reservations for res %s", res.Id)
				s.app.Event.Emit("sync:reservations", "updated")
			} else {
				log.Printf("ERROR: s.app is nil, cannot emit event")
			}
		}
	}
}

func (s *RiadService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *RiadService) GetLocalRooms() ([]map[string]interface{}, error) {
	return db.GetRooms()
}

func (s *RiadService) GetLocalReservations() ([]map[string]interface{}, error) {
	return db.GetReservations()
}

func (s *RiadService) CreateLocalReservation(userID, roomID, start, end string, amount float64) (string, error) {
	if s.grpcClient != nil {
		resp, err := s.grpcClient.CreateReservation(s.grpcCtx(), &pb.CreateReservationRequest{
			UserId:    userID,
			RoomId:    roomID,
			StartDate: start,
			EndDate:   end,
			Amount:    amount,
		})
		if err == nil && resp.Id != "" {
			if debugLog != nil {
				debugLog.Printf("Reservation created via gRPC: id=%s\n", resp.Id)
			}
			db.SaveReservation(resp.Id, userID, roomID, start, end, amount, resp.Status)
			db.MarkSynced("reservations", resp.Id)
			return resp.Id, nil
		}
		if debugLog != nil {
			debugLog.Printf("gRPC CreateReservation failed: %v. Falling back to local.\n", err)
		}
	}

	rawRooms, err := db.GetRooms()
	if err != nil {
		return "", fmt.Errorf("failed to fetch rooms for validation: %v", err)
	}

	var rooms []logic.Room
	for _, r := range rawRooms {
		rooms = append(rooms, logic.Room{
			ID:          r["id"].(string),
			Number:      r["numero"].(int),
			Type:        r["type"].(string),
			Price:       r["prix"].(float64),
			Description: r["description"].(string),
			Equipments:  r["equipements"].(string),
			Status:      r["statut"].(string),
		})
	}

	res := logic.Reservation{
		ID:        uuid.New().String(),
		UserID:    userID,
		RoomID:    roomID,
		Amount:    amount,
		Status:    "pending",
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	if err := logic.ValidateReservation(res, rooms); err != nil {
		return "", err
	}

	err = db.SaveReservation(res.ID, res.UserID, res.RoomID, start, end, res.Amount, res.Status)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func (s *RiadService) UpdateLocalRoom(id string, num int, roomType string, price float64, desc, equip, status string) error {
	if debugLog != nil {
		debugLog.Printf("UpdateLocalRoom called: ID=%s, Num=%d, Price=%.2f\n", id, num, price)
	}
	room := logic.Room{
		ID:          id,
		Number:      num,
		Type:        roomType,
		Price:       price,
		Description: desc,
		Equipments:  equip,
		Status:      status,
	}

	if err := logic.ValidateRoom(room); err != nil {
		if debugLog != nil {
			debugLog.Printf("Validation failed for room %s: %v\n", id, err)
		}
		return err
	}

	err := db.SaveRoom(id, num, roomType, price, desc, equip, status, "propre")
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("db.SaveRoom failed for room %s: %v\n", id, err)
		}
		return err
	}
	if debugLog != nil {
		debugLog.Printf("Room %s saved successfully to local DB\n", id)
	}
	return nil
}

func (s *RiadService) UpdateLocalReservation(id string, userId, roomId, start, end string, amount float64, status string) error {
	if debugLog != nil {
		debugLog.Printf("UpdateLocalReservation called: ID=%s\n", id)
	}
	err := db.SaveReservation(id, userId, roomId, start, end, amount, status)
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("db.SaveReservation failed for reservation %s: %v\n", id, err)
		}
		return err
	}
	if debugLog != nil {
		debugLog.Printf("Reservation %s updated successfully in local DB\n", id)
	}
	return nil
}

func (s *RiadService) UpdateCleaningStatus(id, status string) error {
	if debugLog != nil {
		debugLog.Printf("UpdateCleaningStatus called: ID=%s, Status=%s\n", id, status)
	}

	if s.grpcClient != nil {
		_, err := s.grpcClient.UpdateCleaningStatus(s.grpcCtx(), &pb.UpdateCleaningStatusRequest{
			RoomId:         id,
			CleaningStatus: status,
		})
		if err == nil {
			if debugLog != nil {
				debugLog.Printf("Cleaning status updated via gRPC for room %s\n", id)
			}
			if s.app != nil {
				s.app.Event.Emit("sync:rooms", "updated")
			}
			return nil
		}
		if debugLog != nil {
			debugLog.Printf("gRPC UpdateCleaningStatus failed: %v. Falling back to local.\n", err)
		}
	}

	rawRooms, err := db.GetRooms()
	if err != nil {
		return err
	}

	var room map[string]interface{}
	for _, r := range rawRooms {
		if r["id"] == id {
			room = r
			break
		}
	}

	if room == nil {
		return fmt.Errorf("room not found")
	}

	err = db.SaveRoom(
		id,
		int(room["numero"].(int)),
		room["type"].(string),
		room["prix"].(float64),
		room["description"].(string),
		room["equipements"].(string),
		room["statut"].(string),
		status,
	)

	if err != nil {
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
	return db.GetUnsynced(table)
}

func (s *RiadService) StartSyncLoop() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		if debugLog != nil {
			debugLog.Println("Background sync loop started (Debug mode: 10s)")
		}

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
		if debugLog != nil {
			debugLog.Println("Sync skipped: no token present")
		}
		return
	}
	if debugLog != nil {
		debugLog.Println("Running gRPC background sync...")
	}

	s.gRPCPullUpdates()
	s.gRPCSyncReservations()
}

func (s *RiadService) gRPCPullUpdates() {
	if s.grpcClient == nil {
		if debugLog != nil {
			debugLog.Println("gRPC pull skipped: no client")
		}
		return
	}

	resp, err := s.grpcClient.SyncData(s.grpcCtx(), &pb.SyncDataRequest{
		LastSequenceId: s.lastSyncTimestamp,
	})
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("gRPC SyncData failed: %v\n", err)
		}
		return
	}

	for _, room := range resp.Rooms {
		if err := db.SaveRoom(room.Id, int(room.Number), room.Type, room.Price, room.Description, room.Equipments, room.Status, room.CleaningStatus); err != nil {
			if debugLog != nil {
				debugLog.Printf("Error saving pulled room %s: %v\n", room.Id, err)
			}
		}
	}

	for _, res := range resp.Reservations {
		if err := db.SaveReservation(res.Id, res.UserId, res.RoomId, res.StartDate, res.EndDate, res.Amount, res.Status); err != nil {
			if debugLog != nil {
				debugLog.Printf("Error saving pulled reservation %s: %v\n", res.Id, err)
			}
		}
		db.MarkSynced("reservations", res.Id)
	}

	s.lastSyncTimestamp = time.Now().Unix()
	if debugLog != nil {
		debugLog.Printf("Successfully synced %d rooms and %d reservations from server via gRPC\n", len(resp.Rooms), len(resp.Reservations))
	}
}

func (s *RiadService) gRPCSyncReservations() {
	if s.grpcClient == nil {
		return
	}

	unsynced, err := db.GetUnsynced("reservations")
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("Error fetching unsynced reservations: %v\n", err)
		}
		return
	}

	if len(unsynced) == 0 {
		return
	}

	for _, resMap := range unsynced {
		userID, _ := resMap["user_id"].(string)
		roomID, _ := resMap["chambre_id"].(string)
		start, _ := resMap["date_debut"].(string)
		end, _ := resMap["date_fin"].(string)
		amount, _ := resMap["montant"].(float64)
		id, _ := resMap["id"].(string)

		resp, err := s.grpcClient.CreateReservation(s.grpcCtx(), &pb.CreateReservationRequest{
			UserId:    userID,
			RoomId:    roomID,
			StartDate: start,
			EndDate:   end,
			Amount:    amount,
		})
		if err != nil {
			if debugLog != nil {
				debugLog.Printf("gRPC sync failed for reservation %s: %v\n", id, err)
			}
			continue
		}

		if resp.Id != "" {
			if err := db.MarkSynced("reservations", id); err != nil {
				if debugLog != nil {
					debugLog.Printf("Error marking reservation %s as synced: %v\n", id, err)
				}
			} else if debugLog != nil {
				debugLog.Printf("Successfully synced reservation %s via gRPC (server id: %s)\n", id, resp.Id)
			}
		}
	}
}
