package logic

import (
	"errors"

	"RIAD_SERVER/internal/eventbus"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"gorm.io/gorm"
)

func ValidateReservation(r Reservation, c Chambre) error {
	if r.DateDebut >= r.DateFin {
		return errors.New("dates invalides")
	}
	if !c.CanBook() {
		return errors.New("chambre non disponible")
	}
	return nil
}

func (r *Reservation) Checkin(c *Chambre) error {
	if r.Statut != "confirmée" {
		return errors.New("réservation non confirmée")
	}
	r.Statut = "checkin"
	c.Statut = "occupee"
	return nil
}

func (r *Reservation) Checkout(c *Chambre) error {
	if r.Statut != "checkin" {
		return errors.New("réservation n'est pas en checkin")
	}
	r.Statut = "checkout"
	c.Statut = "libre"
	return nil
}

func GetReservationsByUser(db *gorm.DB, userID string) ([]Reservation, error) {
	var reservations []Reservation
	if err := db.Where("user_id = ?", userID).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservations(db *gorm.DB) ([]Reservation, error) {
	var reservations []Reservation
	if err := db.Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func CreateReservation(db *gorm.DB, res *Reservation) error {
	var chambre Chambre
	if err := db.First(&chambre, "id = ?", res.ChambreID).Error; err != nil {
		return errors.New("chambre non trouvée")
	}

	if err := ValidateReservation(*res, chambre); err != nil {
		return err
	}

	if err := db.Create(res).Error; err != nil {
		return errors.New("erreur lors de la création de la réservation")
	}

	eventbus.GlobalBus.Publish("RESERVATION_UPDATED", res.ID, &pb.Reservation{
		Id:        res.ID,
		UserId:    res.UserID,
		RoomId:    res.ChambreID,
		StartDate: res.DateDebut,
		EndDate:   res.DateFin,
		Amount:    res.Montant,
		Status:    res.Statut,
	})

	return nil
}

func CheckinReservation(db *gorm.DB, id string) (*Reservation, error) {
	var res Reservation
	if err := db.First(&res, "id = ?", id).Error; err != nil {
		return nil, errors.New("réservation non trouvée")
	}

	var chambre Chambre
	if err := db.First(&chambre, "id = ?", res.ChambreID).Error; err != nil {
		return nil, errors.New("chambre non trouvée")
	}

	if err := res.Checkin(&chambre); err != nil {
		return nil, err
	}

	if err := db.Save(&res).Error; err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la réservation")
	}
	if err := db.Save(&chambre).Error; err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la chambre")
	}

	eventbus.GlobalBus.Publish("RESERVATION_UPDATED", res.ID, &pb.Reservation{
		Id:        res.ID,
		UserId:    res.UserID,
		RoomId:    res.ChambreID,
		StartDate: res.DateDebut,
		EndDate:   res.DateFin,
		Amount:    res.Montant,
		Status:    res.Statut,
	})
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

	return &res, nil
}

func CheckoutReservation(db *gorm.DB, id string) (*Reservation, error) {
	var res Reservation
	if err := db.First(&res, "id = ?", id).Error; err != nil {
		return nil, errors.New("réservation non trouvée")
	}

	var chambre Chambre
	if err := db.First(&chambre, "id = ?", res.ChambreID).Error; err != nil {
		return nil, errors.New("chambre non trouvée")
	}

	if err := res.Checkout(&chambre); err != nil {
		return nil, err
	}

	if err := db.Save(&res).Error; err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la réservation")
	}
	if err := db.Save(&chambre).Error; err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la chambre")
	}

	eventbus.GlobalBus.Publish("RESERVATION_UPDATED", res.ID, &pb.Reservation{
		Id:        res.ID,
		UserId:    res.UserID,
		RoomId:    res.ChambreID,
		StartDate: res.DateDebut,
		EndDate:   res.DateFin,
		Amount:    res.Montant,
		Status:    res.Statut,
	})
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

	return &res, nil
}