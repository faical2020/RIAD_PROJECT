package logic

import (
	"gorm.io/gorm"
)

type SyncResponse struct {
	Chambres     []Chambre     `json:"chambres"`
	Reservations []Reservation `json:"reservations"`
}

func GetSyncUpdates(db *gorm.DB, since int64) (*SyncResponse, error) {
	var chambres []Chambre
	if err := db.Where("updated_at > ?", since).Find(&chambres).Error; err != nil {
		return nil, err
	}

	var reservations []Reservation
	if err := db.Where("updated_at > ?", since).Find(&reservations).Error; err != nil {
		return nil, err
	}

	return &SyncResponse{
		Chambres:     chambres,
		Reservations: reservations,
	}, nil
}
