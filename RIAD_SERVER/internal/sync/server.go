package sync

import (
	"log"

	pb "github.com/anomalyco/riad_project/proto/sync"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SyncServer struct {
	pb.UnimplementedSyncServiceServer
}

func NewSyncServer() *SyncServer {
	return &SyncServer{}
}

func (s *SyncServer) StreamUpdates(req *pb.SyncRequest, stream pb.SyncService_StreamUpdatesServer) error {
	// 1. Basic Authentication
	if req.Token == "" {
		return status.Error(codes.Unauthenticated, "token is required")
	}
	
	log.Printf("Client connected to sync stream. Token: %s...", req.Token[:5])

	// 2. Subscribe to the Global Event Bus
	eventChan := GlobalBus.Subscribe()
	defer func() {
		// Note: In a production system, we would need a way to unsubscribe
		log.Printf("Client disconnected from sync stream")
	}()

	// 3. Stream events from the bus to the gRPC stream
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
		case "ROOM_DELETED":
			syncEvent.Type = pb.SyncEvent_ROOM_DELETED
		case "RESERVATION_DELETED":
			syncEvent.Type = pb.SyncEvent_RESERVATION_DELETED
		default:
			log.Printf("Unknown event type: %s", event.Type)
			continue
		}

		if err := stream.Send(&syncEvent); err != nil {
			log.Printf("Error sending event to stream: %v", err)
			return err
		}
	}

	return nil
}
