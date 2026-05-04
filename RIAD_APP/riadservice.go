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
	ctx context.Context
}

func NewRiadService() *RiadService {
	return &RiadService{}
}

func (s *RiadService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *RiadService) GetLocalRooms() ([]map[string]interface{}, error) {
	return db.GetRooms()
}

func (s *RiadService) CreateLocalReservation(userID, roomID, start, end string, amount float64) (string, error) {
	res := logic.Reservation{
		ID:        uuid.New().String(),
		UserID:    userID,
		RoomID:    roomID,
		Amount:    amount,
		Status:    "pending",
		StartDate: time.Now(), // Simplified for now, should parse 'start'
		EndDate:   time.Now(), // Simplified
	}

	if err := logic.ValidateReservation(res, nil); err != nil {
		return "", err
	}

	err := db.SaveReservation(res.ID, res.UserID, res.RoomID, start, end, res.Amount, res.Status)
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

	err := db.SaveRoom(id, num, roomType, price, desc, equip, status)
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

func (s *RiadService) MarkAsSynced(table, id string) error {
	return db.MarkSynced(table, id)
}

func (s *RiadService) GetUnsynced(table string) ([]map[string]interface{}, error) {
	return db.GetUnsynced(table)
}
