package db

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&Room{}, &Reservation{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestSaveAndGetRoom(t *testing.T) {
	db := setupTestDB(t)

	db.Save(&Room{ID: "test-id", Numero: 101, Type: "double", Prix: 500})

	var room Room
	if err := db.First(&room, "id = ?", "test-id").Error; err != nil {
		t.Fatalf("GetRoom failed: %v", err)
	}
	if room.Numero != 101 {
		t.Errorf("expected numero 101, got %d", room.Numero)
	}
	if room.Type != "double" {
		t.Errorf("expected type double, got %s", room.Type)
	}
}

func TestSaveAndGetReservation(t *testing.T) {
	db := setupTestDB(t)

	res := Reservation{
		ID: "res-test-1", UserID: "user-1", ChambreID: "room-1",
		DateDebut: "2026-06-01", DateFin: "2026-06-05",
		Montant: 2000, Statut: "confirmée",
	}
	db.Save(&res)

	var fetched Reservation
	db.First(&fetched, "id = ?", "res-test-1")
	if fetched.Montant != 2000 {
		t.Errorf("expected 2000, got %f", fetched.Montant)
	}
}

func TestMarkSynced(t *testing.T) {
	db := setupTestDB(t)

	db.Save(&Reservation{ID: "res-sync", Statut: "pending", Synced: false})
	db.Table("reservations").Where("id = ?", "res-sync").Update("synced", true)

	var r Reservation
	db.First(&r, "id = ?", "res-sync")
	if !r.Synced {
		t.Error("expected synced to be true")
	}
}

func TestGetUnsynced(t *testing.T) {
	db := setupTestDB(t)

	db.Where("1 = 1").Delete(&Reservation{})
	db.Save(&Reservation{ID: "unsynced-1", Statut: "pending", Synced: false})
	db.Save(&Reservation{ID: "unsynced-2", Statut: "confirmed", Synced: true})

	var unsynced []Reservation
	db.Where("synced = ?", false).Find(&unsynced)

	if len(unsynced) != 1 {
		t.Errorf("expected 1 unsynced, got %d", len(unsynced))
	}
}

func TestRooms(t *testing.T) {
	db := setupTestDB(t)

	db.Save(&Room{ID: "r1", Numero: 1, Type: "single"})
	db.Save(&Room{ID: "r2", Numero: 2, Type: "double"})

	var rooms []Room
	db.Find(&rooms)
	if len(rooms) != 2 {
		t.Errorf("expected 2 rooms, got %d", len(rooms))
	}
}

func TestReservations(t *testing.T) {
	db := setupTestDB(t)

	db.Save(&Reservation{ID: "resa1", UserID: "u1", ChambreID: "c1", Statut: "confirmée"})
	db.Save(&Reservation{ID: "resa2", UserID: "u2", ChambreID: "c2", Statut: "en attente"})

	var reservations []Reservation
	db.Find(&reservations)
	if len(reservations) != 2 {
		t.Errorf("expected 2 reservations, got %d", len(reservations))
	}
}
