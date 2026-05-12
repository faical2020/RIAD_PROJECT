package tests

import (
	"os"
	"testing"
	"time"

	"RIAD_APP/internal/db"
	"RIAD_APP/pkg/logic"
)

func TestHybridLogic(t *testing.T) {
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

		invalidRoom := logic.Room{Number: 0, Price: 1500}
		if err := logic.ValidateRoom(invalidRoom); err == nil {
			t.Error("Expected invalid room number (0) to fail")
		}

		invalidPrice := logic.Room{Number: 101, Price: -100}
		if err := logic.ValidateRoom(invalidPrice); err == nil {
			t.Error("Expected negative price to fail")
		}
	})

	t.Run("ValidateReservation", func(t *testing.T) {
		start := time.Now()
		end := start.Add(48 * time.Hour)

		validRes := logic.Reservation{
			UserID:    "user-1",
			RoomID:    "room-1",
			StartDate: start,
			EndDate:   end,
			Amount:    1500,
		}
		if err := logic.ValidateReservation(validRes, nil); err != nil {
			t.Errorf("Expected valid reservation to pass, got: %v", err)
		}

		invalidDates := logic.Reservation{
			UserID:    "user-1",
			RoomID:    "room-1",
			StartDate: end,
			EndDate:   start,
			Amount:    1500,
		}
		if err := logic.ValidateReservation(invalidDates, nil); err == nil {
			t.Error("Expected reservation with inverted dates to fail")
		}

		missingUser := logic.Reservation{
			UserID:    "",
			RoomID:    "room-1",
			StartDate: start,
			EndDate:   end,
			Amount:    1500,
		}
		if err := logic.ValidateReservation(missingUser, nil); err == nil {
			t.Error("Expected reservation with empty UserID to fail")
		}

		negativeAmount := logic.Reservation{
			UserID:    "user-1",
			RoomID:    "room-1",
			StartDate: start,
			EndDate:   end,
			Amount:    -500,
		}
		if err := logic.ValidateReservation(negativeAmount, nil); err == nil {
			t.Error("Expected reservation with negative amount to fail")
		}
	})

	t.Run("LocalPersistence", func(t *testing.T) {
		roomID := "room-1"
		err := db.SaveRoom(roomID, 101, "Double", 1500.0, "Desc", "Wifi", "libre", "propre")
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
