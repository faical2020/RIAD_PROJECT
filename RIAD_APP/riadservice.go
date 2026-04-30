package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"RIAD_APP/internal/db"
	"RIAD_APP/pkg/logic"
)

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
		return err
	}

	return db.SaveRoom(id, num, roomType, price, desc, equip, status)
}

func (s *RiadService) MarkAsSynced(table, id string) error {
	return db.MarkSynced(table, id)
}

func (s *RiadService) GetUnsynced(table string) ([]map[string]interface{}, error) {
	return db.GetUnsynced(table)
}
