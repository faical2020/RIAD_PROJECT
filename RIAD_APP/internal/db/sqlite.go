package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Room struct {
	ID             string  `json:"id" gorm:"primaryKey"`
	Numero         int     `json:"numero"`
	Type           string  `json:"type"`
	Prix           float64 `json:"prix"`
	Description    string  `json:"description"`
	Equipements    string  `json:"equipements"`
	Statut         string  `json:"statut"`
	CleaningStatus string  `json:"cleaning_status" gorm:"column:cleaning_status;default:propre"`
	Synced         bool    `json:"synced" gorm:"default:1"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now().Unix()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

func (r *Room) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now().Unix()
	return nil
}

type Reservation struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	UserID    string  `json:"user_id" gorm:"column:user_id"`
	ChambreID string  `json:"chambre_id" gorm:"column:chambre_id"`
	DateDebut string  `json:"date_debut" gorm:"column:date_debut"`
	DateFin   string  `json:"date_fin" gorm:"column:date_fin"`
	Montant   float64 `json:"montant"`
	Statut    string  `json:"statut"`
	Synced    bool    `json:"synced" gorm:"default:0"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	now := time.Now().Unix()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

func (r *Reservation) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now().Unix()
	return nil
}

func GetDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "riad.db"
	}
	return filepath.Join(configDir, "riad_app", "riad.db")
}

func InitDB(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erreur creation dossier db: %v", err)
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("erreur ouverture db: %v", err)
	}

	if err := DB.AutoMigrate(&Room{}, &Reservation{}); err != nil {
		return fmt.Errorf("erreur migration: %v", err)
	}
	return nil
}

func GetRooms() ([]Room, error) {
	var rooms []Room
	if err := DB.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func GetRoomByID(id string) (*Room, error) {
	var room Room
	if err := DB.First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func GetReservations() ([]Reservation, error) {
	var reservations []Reservation
	if err := DB.Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func SaveRoom(id string, num int, roomType string, price float64, desc, equip, status, cleaningStatus string) error {
	room := Room{
		ID:             id,
		Numero:         num,
		Type:           roomType,
		Prix:           price,
		Description:    desc,
		Equipements:    equip,
		Statut:         status,
		CleaningStatus: cleaningStatus,
	}
	return DB.Save(&room).Error
}

func SaveReservation(id, userID, chambreID, dateDebut, dateFin string, montant float64, statut string) error {
	res := Reservation{
		ID:        id,
		UserID:    userID,
		ChambreID: chambreID,
		DateDebut: dateDebut,
		DateFin:   dateFin,
		Montant:   montant,
		Statut:    statut,
	}
	return DB.Save(&res).Error
}

func MarkSynced(table string, id string) error {
	return DB.Table(table).Where("id = ?", id).Update("synced", true).Error
}

func GetUnsynced(table string) ([]Reservation, error) {
	var reservations []Reservation
	if err := DB.Table(table).Where("synced = ?", false).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}
