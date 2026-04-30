package tests

import (
	"testing"
	"RIAD_APP/internal/db"
	"RIAD_APP/pkg/logic"
	"os"
)

func TestHybridLogic(t *testing.T) {
	// Setup temporary DB
	dbPath := "test_riad.db"
	err := db.InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer os.Remove(dbPath)

	t.Run("RoomValidation", func(t *testing.T) {
		validRoom := logic.Room{Number: 101, Price: 1500}
		if err := logic.ValidateRoom(validRoom); err != nil {
			t.Errorf("Expected valid room to pass, got: %v", err)
		}

		invalidRoom := logic.Room{Number: -1, Price: 1500}
		if err := logic.ValidateRoom(invalidRoom); err == nil {
			t.Error("Expected invalid room number to fail")
		}
	})

	t.Run("LocalPersistence", func(t *testing.T) {
		roomID := "room-1"
		err := db.SaveRoom(roomID, 101, "Double", 1500.0, "Desc", "Wifi", "libre")
		if err != nil {
			t.Fatalf("Failed to save room: %v", err)
		}

		rooms, err := db.GetRooms()
		if err != nil {
			t.Fatalf("Failed to get rooms: %v", err)
		}

		found := false
		for _, r := range rooms {
			if r["id"] == roomID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Room was not found in local SQLite")
		}
	})

	t.Run("SyncFlagging", func(t *testing.T) {
		resID := "res-1"
		err := db.SaveReservation(resID, "user-1", "room-1", "2026-01-01", "2026-01-05", 5000.0, "pending")
		if err != nil {
			t.Fatalf("Failed to save reservation: %v", err)
		}

		unsynced, err := db.GetUnsynced("reservations")
		if err != nil {
			t.Fatalf("Failed to get unsynced: %v", err)
		}

		if len(unsynced) == 0 {
			t.Error("Reservation should be marked as unsynced by default")
		}

		err = db.MarkSynced("reservations", resID)
		if err != nil {
			t.Fatalf("Failed to mark synced: %v", err)
		}

		unsynced, _ = db.GetUnsynced("reservations")
		if len(unsynced) != 0 {
			t.Error("Reservation should no longer be unsynced")
		}
	})
}
