package logic

import (
	"errors"
	"time"

	"RIAD_SERVER/internal/eventbus"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"gorm.io/gorm"
)

func ValidateChambre(c Chambre) error {
	if c.Numero <= 0 {
		return errors.New("numéro de chambre invalide")
	}
	if c.Prix <= 0 {
		return errors.New("prix invalide")
	}
	return nil
}

func (c *Chambre) CanBook() bool {
	return c.Statut == "libre"
}

func GetChambres(db *gorm.DB) ([]Chambre, error) {
	var chambres []Chambre
	if err := db.Find(&chambres).Error; err != nil {
		return nil, err
	}
	return chambres, nil
}

func UpdateCleaningStatus(db *gorm.DB, id, status string) (*Chambre, error) {
	var chambre Chambre
	if err := db.First(&chambre, "id = ?", id).Error; err != nil {
		return nil, errors.New("chambre non trouvée")
	}

	chambre.CleaningStatus = status
	chambre.UpdatedAt = time.Now().Unix()
	if err := db.Save(&chambre).Error; err != nil {
		return nil, errors.New("erreur lors de la mise à jour")
	}

	eventbus.GlobalBus.Publish("ROOM_UPDATED", chambre.ID, &pb.Room{
		Id:             chambre.ID,
		Number:         int32(chambre.Numero),
		Type:           chambre.Type,
		Price:          chambre.Prix,
		Description:    chambre.Description,
		Equipments:     chambre.Equipements,
		Status:         chambre.Statut,
		CleaningStatus: chambre.CleaningStatus,
	})

	return &chambre, nil
}