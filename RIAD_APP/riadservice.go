package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"RIAD_APP/internal/db"
	"RIAD_APP/pkg/logic"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"github.com/wailsapp/wails/v3/pkg/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	ctx                context.Context
	token              string
	apiBaseURL         string
	app                *application.App // Reference to the Wails app for emitting events
	lastSyncTimestamp  int64
}

func NewRiadService() *RiadService {
	return &RiadService{
		apiBaseURL: "http://localhost:8081/api/v1",
	}
}

func (s *RiadService) SetApp(app *application.App) {
	s.app = app
}

func (s *RiadService) SetToken(token string) {
	s.token = token
	if debugLog != nil {
		debugLog.Printf("SetToken called: token received (length: %d)\n", len(token))
	}
	// Trigger immediate sync when token is set
	go s.performSync()
	// Start the real-time gRPC stream
	go s.startGRPCStream()
}

func (s *RiadService) startGRPCStream() {
	if s.token == "" {
		return
	}

	// Connection settings
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("gRPC Dial failed: %v\n", err)
		}
		return
	}
	defer conn.Close()

	client := pb.NewSyncServiceClient(conn)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if debugLog != nil {
				debugLog.Println("Attempting to connect to gRPC sync stream...")
			}

			stream, err := client.StreamUpdates(s.ctx, &pb.SyncRequest{
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
					break // Break inner loop to reconnect
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
			// Note: CleaningStatus is not yet in the proto Room message, using default 'propre'
			db.SaveRoom(room.Id, int(room.Number), room.Type, room.Price, room.Description, room.Equipments, room.Status, "propre")
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
	// Fetch the rooms for validation from local DB
	rawRooms, err := db.GetRooms()
	if err != nil {
		return "", fmt.Errorf("failed to fetch rooms for validation: %v", err)
	}

	// Convert []map[string]interface{} to []logic.Room as expected by ValidateReservation
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

	// Pass the slice of rooms
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
	
	// Update local DB
	// We need to fetch the room first to keep other fields
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
	
	// In SQLite we use SaveRoom, but SaveRoom replaces the whole row.
	// For simplicity and consistency with our current SQLite implementation:
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
	
	// Trigger sync event to frontend
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

		// Initial sync attempt
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
		debugLog.Println("Running optimized background sync...")
	}

	// 1. Pull updated data from server using the /sync endpoint
	s.pullUpdates()

	// 2. Push unsynced local data to server
	s.syncReservations()
}

func (s *RiadService) pullUpdates() {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/sync?since=%d", s.apiBaseURL, s.lastSyncTimestamp)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := client.Do(req)
	if err != nil {
		if debugLog != nil {
			debugLog.Printf("Error pulling updates: %v\n", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var syncData struct {
			Chambres     []map[string]interface{} `json:"chambres"`
			Reservations []map[string]interface{} `json:"reservations"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&syncData); err != nil {
			if debugLog != nil {
				debugLog.Printf("Error decoding sync data: %v\n", err)
			}
			return
		}

		// Update local rooms
		for _, r := range syncData.Chambres {
			id, _ := r["id"].(string)
			num, _ := r["numero"].(float64)
			roomType, _ := r["type"].(string)
			price, _ := r["prix"].(float64)
			desc, _ := r["description"].(string)
			equip, _ := r["equipements"].(string)
			status, _ := r["statut"].(string)
			cleaningStatus, _ := r["cleaning_status"].(string)
			if cleaningStatus == "" {
				cleaningStatus = "propre"
			}

			if err := db.SaveRoom(id, int(num), roomType, price, desc, equip, status, cleaningStatus); err != nil {
				if debugLog != nil {
					debugLog.Printf("Error saving pulled room %s: %v\n", id, err)
				}
			}
		}

		// Update local reservations
		for _, r := range syncData.Reservations {
			id, _ := r["id"].(string)
			uId, _ := r["user_id"].(string)
			rId, _ := r["chambre_id"].(string)
			start, _ := r["date_debut"].(string)
			end, _ := r["date_fin"].(string)
			amount, _ := r["montant"].(float64)
			status, _ := r["statut"].(string)

			if err := db.SaveReservation(id, uId, rId, start, end, amount, status); err != nil {
				if debugLog != nil {
					debugLog.Printf("Error saving pulled reservation %s: %v\n", id, err)
				}
			}
			db.MarkSynced("reservations", id)
		}

		s.lastSyncTimestamp = time.Now().Unix()
		if debugLog != nil {
			debugLog.Printf("Successfully synced %d rooms and %d reservations from server\n", len(syncData.Chambres), len(syncData.Reservations))
		}
	}
}

func (s *RiadService) syncReservations() {
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

	client := &http.Client{Timeout: 10 * time.Second}

	for _, resMap := range unsynced {
		// Map map[string]interface{} to the expected JSON payload
		payload := map[string]interface{}{
			"user_id":    resMap["user_id"],
			"chambre_id": resMap["chambre_id"],
			"date_debut": resMap["date_debut"],
			"date_fin":   resMap["date_fin"],
			"montant":    resMap["montant"],
		}

		jsonPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", s.apiBaseURL+"/reservations", bytes.NewBuffer(jsonPayload))
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := client.Do(req)
		if err != nil {
			if debugLog != nil {
				debugLog.Printf("Sync failed for reservation %v: %v\n", resMap["id"], err)
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			id := fmt.Sprintf("%v", resMap["id"])
			if err := db.MarkSynced("reservations", id); err != nil {
				if debugLog != nil {
					debugLog.Printf("Error marking reservation %s as synced: %v\n", id, err)
				}
			} else if debugLog != nil {
				debugLog.Printf("Successfully synced reservation %s\n", id)
			}
		} else {
			body, _ := io.ReadAll(resp.Body)
			if debugLog != nil {
				debugLog.Printf("Server rejected reservation %v with status %d: %s\n", resMap["id"], resp.StatusCode, string(body))
			}
		}
	}
}
