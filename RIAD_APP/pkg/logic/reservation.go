package logic

import (
	"errors"
	"time"
)

type Reservation struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RoomID     string    `json:"chambre_id"`
	StartDate  time.Time `json:"date_debut"`
	EndDate    time.Time `json:"date_fin"`
	Amount     float64   `json:"montant"`
	Status     string    `json:"statut"`
	Synced     bool      `json:"synced"`
}

func ValidateReservation(res Reservation, existingRooms []Room) error {
	if res.UserID == "" || res.RoomID == "" {
		return errors.New("l'utilisateur et la chambre sont requis")
	}
	if res.EndDate.Before(res.StartDate) {
		return errors.New("la date de fin doit être après la date de début")
	}
	if res.Amount <= 0 {
		return errors.New("le montant doit être supérieur à zéro")
	}
	return nil
}
