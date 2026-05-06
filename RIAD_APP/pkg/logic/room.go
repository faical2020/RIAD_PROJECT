package logic

import "errors"

type Room struct {
	ID             string  `json:"id"`
	Number         int     `json:"numero"`
	Type           string  `json:"type"`
	Price          float64 `json:"prix"`
	Description    string  `json:"description"`
	Equipments     string  `json:"equipements"`
	Status         string  `json:"statut"`
	CleaningStatus string  `json:"cleaning_status"`
}

func ValidateRoom(room Room) error {
	if room.Number <= 0 {
		return errors.New("le numéro de chambre doit être positif")
	}
	if room.Price < 0 {
		return errors.New("le prix ne peut pas être négatif")
	}
	return nil
}
