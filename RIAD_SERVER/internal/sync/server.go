package sync

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/eventbus"
	"RIAD_SERVER/internal/logic"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SyncServer struct {
	pb.UnimplementedSyncServiceServer
}

func NewSyncServer() *SyncServer {
	return &SyncServer{}
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		vals = md.Get("token")
	}
	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing token")
	}
	return strings.TrimPrefix(vals[0], "Bearer "), nil
}

func (s *SyncServer) StreamUpdates(req *pb.SyncRequest, stream pb.SyncService_StreamUpdatesServer) error {
	if req.Token == "" {
		return status.Error(codes.Unauthenticated, "token is required")
	}

	slog.Info("client connected to sync stream", "token_prefix", req.Token[:5])

	eventChan := eventbus.GlobalBus.Subscribe()
	defer func() {
		slog.Info("client disconnected from sync stream")
	}()

	for event := range eventChan {
		var syncEvent pb.SyncEvent
		syncEvent.EntityId = event.EntityID
		syncEvent.SequenceId = event.SequenceID

		switch event.Type {
		case "ROOM_UPDATED":
			syncEvent.Type = pb.SyncEvent_ROOM_UPDATED
			if r, ok := event.Data.(*pb.Room); ok {
				syncEvent.Data = &pb.SyncEvent_Room{Room: r}
			}
		case "RESERVATION_UPDATED":
			syncEvent.Type = pb.SyncEvent_RESERVATION_UPDATED
			if res, ok := event.Data.(*pb.Reservation); ok {
				syncEvent.Data = &pb.SyncEvent_Reservation{Reservation: res}
			}
		case "CONSOMMATION_UPDATED":
			syncEvent.Type = pb.SyncEvent_CONSOMMATION_UPDATED
			if c, ok := event.Data.(map[string]interface{}); ok {
				conso := &pb.Consommation{
					Id:            toString(c["id"]),
					ReservationId: toString(c["reservation_id"]),
					Libelle:       toString(c["libelle"]),
					Quantite:      toInt32(c["quantite"]),
					PrixUnitaire:  toFloat64(c["prix_unitaire"]),
				}
				syncEvent.Data = &pb.SyncEvent_Consommation{Consommation: conso}
			}
		case "ROOM_DELETED":
			syncEvent.Type = pb.SyncEvent_ROOM_DELETED
		case "RESERVATION_DELETED":
			syncEvent.Type = pb.SyncEvent_RESERVATION_DELETED
		default:
			slog.Warn("unknown event type", "type", event.Type)
			continue
		}

		if err := stream.Send(&syncEvent); err != nil {
			slog.Error("error sending event to stream", "error", err)
			return err
		}
	}

	return nil
}

func (s *SyncServer) GetChambres(ctx context.Context, req *pb.GetChambresRequest) (*pb.GetChambresResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	chambres, err := logic.GetChambres(db.GetDB())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rooms := make([]*pb.Room, len(chambres))
	for i, c := range chambres {
		rooms[i] = &pb.Room{
			Id: c.ID, Number: int32(c.Numero), Type: c.Type, Price: c.Prix,
			Description: c.Description, Equipments: c.Equipements,
			Status: c.Statut, CleaningStatus: c.CleaningStatus,
		}
	}

	return &pb.GetChambresResponse{Rooms: rooms}, nil
}

func (s *SyncServer) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	res := logic.Reservation{
		UserID: req.UserId, ChambreID: req.RoomId,
		DateDebut: req.StartDate, DateFin: req.EndDate, Montant: req.Amount,
	}

	if err := logic.CreateReservation(db.GetDB(), &res); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateReservationResponse{Id: res.ID, Status: res.Statut}, nil
}

func (s *SyncServer) UpdateCleaningStatus(ctx context.Context, req *pb.UpdateCleaningStatusRequest) (*pb.UpdateCleaningStatusResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	chambre, err := logic.UpdateCleaningStatus(db.GetDB(), req.RoomId, req.CleaningStatus)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdateCleaningStatusResponse{
		Room: &pb.Room{
			Id: chambre.ID, Number: int32(chambre.Numero), Type: chambre.Type,
			Price: chambre.Prix, Description: chambre.Description,
			Equipments: chambre.Equipements, Status: chambre.Statut,
			CleaningStatus: chambre.CleaningStatus,
		},
	}, nil
}

func (s *SyncServer) Checkin(ctx context.Context, req *pb.CheckinRequest) (*pb.CheckinResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	res, err := logic.CheckinReservation(db.GetDB(), req.ReservationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CheckinResponse{Reservation: &pb.Reservation{
		Id: res.ID, UserId: res.UserID, RoomId: res.ChambreID,
		StartDate: res.DateDebut, EndDate: res.DateFin,
		Amount: res.Montant, Status: res.Statut,
	}}, nil
}

func (s *SyncServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	res, err := logic.CheckoutReservation(db.GetDB(), req.ReservationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CheckoutResponse{Reservation: &pb.Reservation{
		Id: res.ID, UserId: res.UserID, RoomId: res.ChambreID,
		StartDate: res.DateDebut, EndDate: res.DateFin,
		Amount: res.Montant, Status: res.Statut,
	}}, nil
}

func (s *SyncServer) SyncData(ctx context.Context, req *pb.SyncDataRequest) (*pb.SyncDataResponse, error) {
	if _, _, err := extractTokenAndValidate(ctx); err != nil {
		return nil, err
	}

	data, err := logic.GetSyncUpdates(db.GetDB(), req.LastSequenceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rooms := make([]*pb.Room, len(data.Chambres))
	for i, c := range data.Chambres {
		rooms[i] = &pb.Room{
			Id: c.ID, Number: int32(c.Numero), Type: c.Type, Price: c.Prix,
			Description: c.Description, Equipments: c.Equipements,
			Status: c.Statut, CleaningStatus: c.CleaningStatus,
		}
	}

	reservations := make([]*pb.Reservation, len(data.Reservations))
	for i, r := range data.Reservations {
		reservations[i] = &pb.Reservation{
			Id: r.ID, UserId: r.UserID, RoomId: r.ChambreID,
			StartDate: r.DateDebut, EndDate: r.DateFin,
			Amount: r.Montant, Status: r.Statut,
		}
	}

	return &pb.SyncDataResponse{Rooms: rooms, Reservations: reservations}, nil
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok { return s }
	return fmt.Sprintf("%v", v)
}

func toInt32(v interface{}) int32 {
	switch n := v.(type) {
	case int: return int32(n)
	case int32: return n
	case float64: return int32(n)
	default: return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch n := v.(type) {
	case float64: return n
	case int: return float64(n)
	case int32: return float64(n)
	default: return 0
	}
}

func extractTokenAndValidate(ctx context.Context) (string, string, error) {
	tokenStr, err := extractToken(ctx)
	if err != nil {
		return "", "", err
	}
	userID, role, err := logic.ValidateAccessToken(tokenStr)
	if err != nil {
		return "", "", status.Error(codes.Unauthenticated, "invalid token")
	}
	return userID, role, nil
}
