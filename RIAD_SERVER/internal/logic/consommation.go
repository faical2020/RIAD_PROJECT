package logic

import (
	"errors"
	"time"

	"RIAD_SERVER/internal/eventbus"
	"gorm.io/gorm"
)

// ── Services (catalogue) ──

func GetServices(db *gorm.DB) ([]Service, error) {
	var services []Service
	if err := db.Order("categorie, nom").Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func CreateService(db *gorm.DB, s *Service) error {
	if s.Nom == "" {
		return errors.New("le nom du service est requis")
	}
	return db.Create(s).Error
}

func UpdateService(db *gorm.DB, s *Service) error {
	return db.Save(s).Error
}

func DeleteService(db *gorm.DB, id string) error {
	return db.Delete(&Service{}, "id = ?", id).Error
}

// ── Consommations ──

func GetConsommationsByReservation(db *gorm.DB, reservationID string) ([]Consommation, error) {
	var consommations []Consommation
	if err := db.Where("reservation_id = ?", reservationID).
		Order("ajoute_le ASC").Find(&consommations).Error; err != nil {
		return nil, err
	}
	return consommations, nil
}

func AddConsommation(db *gorm.DB, c *Consommation) error {
	if c.ReservationID == "" || c.Libelle == "" {
		return errors.New("réservation et libellé requis")
	}
	if c.Quantite <= 0 {
		c.Quantite = 1
	}
	if err := db.Create(c).Error; err != nil {
		return err
	}

	eventbus.GlobalBus.Publish("CONSOMMATION_UPDATED", c.ID, map[string]interface{}{
		"id":             c.ID,
		"reservation_id": c.ReservationID,
		"libelle":        c.Libelle,
		"quantite":       c.Quantite,
		"prix_unitaire":  c.PrixUnitaire,
	})

	return nil
}

func DeleteConsommation(db *gorm.DB, id string) error {
	return db.Delete(&Consommation{}, "id = ?", id).Error
}

// ── Facture ──

type FactureDetail struct {
	Reservation       Reservation    `json:"reservation"`
	Chambre           Chambre        `json:"chambre"`
	Consommations     []Consommation `json:"consommations"`
	TotalConsommations float64       `json:"total_consommations"`
	TotalSejour       float64        `json:"total_sejour"`
	Paiements         []Paiement     `json:"paiements"`
	TotalPaye         float64        `json:"total_paye"`
	RestantDu         float64        `json:"restant_du"`
}

func GetFacture(db *gorm.DB, reservationID string) (*FactureDetail, error) {
	var reservation Reservation
	if err := db.First(&reservation, "id = ?", reservationID).Error; err != nil {
		return nil, errors.New("réservation non trouvée")
	}

	var chambre Chambre
	db.First(&chambre, "id = ?", reservation.ChambreID)

	consommations, _ := GetConsommationsByReservation(db, reservationID)

	var totalConso float64
	for _, c := range consommations {
		totalConso += float64(c.Quantite) * c.PrixUnitaire
	}

	var paiements []Paiement
	db.Where("reservation_id = ? AND statut = 'paye'", reservationID).Find(&paiements)

	var totalPaye float64
	for _, p := range paiements {
		totalPaye += p.Montant
	}

	totalSejour := reservation.Montant + totalConso
	restantDu := totalSejour - totalPaye
	if restantDu < 0 {
		restantDu = 0
	}

	return &FactureDetail{
		Reservation:       reservation,
		Chambre:           chambre,
		Consommations:     consommations,
		TotalConsommations: totalConso,
		TotalSejour:       totalSejour,
		Paiements:         paiements,
		TotalPaye:         totalPaye,
		RestantDu:         restantDu,
	}, nil
}

// ── Paiements ──

func AddPaiement(db *gorm.DB, p *Paiement) error {
	if p.ReservationID == "" || p.Montant <= 0 {
		return errors.New("réservation et montant requis")
	}
	if p.ModePaiement == "" {
		p.ModePaiement = "especes"
	}
	p.Statut = "paye"
	p.CreeLe = time.Now().Unix()

	if err := db.Create(p).Error; err != nil {
		return err
	}

	eventbus.GlobalBus.Publish("RESERVATION_UPDATED", p.ReservationID, map[string]string{"id": p.ReservationID})
	return nil
}

func GetPaiementsByReservation(db *gorm.DB, reservationID string) ([]Paiement, error) {
	var paiements []Paiement
	if err := db.Where("reservation_id = ?", reservationID).Find(&paiements).Error; err != nil {
		return nil, err
	}
	return paiements, nil
}
