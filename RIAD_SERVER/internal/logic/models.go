package logic

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        string `json:"id" gorm:"type:uuid;primaryKey"`
    Email     string `json:"email" gorm:"uniqueIndex"`
    Password  string `json:"password,omitempty" gorm:"type:varchar(255)"`
    Nom       string `json:"nom"`
    Prenom    string `json:"prenom"`
    Role      string `json:"role" gorm:"type:varchar(50);default:'client'"`
    Telephone string `json:"telephone"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}

type Chambre struct {
    ID          string `json:"id" gorm:"type:uuid;primaryKey"`
    Numero      int    `json:"numero"`
    Type        string `json:"type" gorm:"type:varchar(50)"`
    Prix        float64 `json:"prix"`
    Statut      string `json:"statut" gorm:"type:varchar(50);default:'libre'"`
    Description string `json:"description"`
    Equipements string `json:"equipements"`
}

func (c *Chambre) BeforeCreate(tx *gorm.DB) error {
    if c.ID == "" {
        c.ID = uuid.New().String()
    }
    return nil
}

type Reservation struct {
    ID         string  `json:"id" gorm:"type:uuid;primaryKey"`
    UserID     string  `json:"user_id" gorm:"type:uuid"`
    ChambreID  string  `json:"chambre_id" gorm:"type:uuid"`
    DateDebut  string  `json:"date_debut"`
    DateFin    string  `json:"date_fin"`
    Statut     string  `json:"statut" gorm:"type:varchar(50);default:'confirmée'"`
    Montant    float64 `json:"montant"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
    if r.ID == "" {
        r.ID = uuid.New().String()
    }
    return nil
}

type Tache struct {
    ID          string `json:"id" gorm:"type:uuid;primaryKey"`
    EmployeID   string `json:"employe_id" gorm:"type:uuid"`
    Description string `json:"description"`
    Statut      string `json:"statut" gorm:"type:varchar(50);default:'à faire'"`
}

func (t *Tache) BeforeCreate(tx *gorm.DB) error {
    if t.ID == "" {
        t.ID = uuid.New().String()
    }
    return nil
}

type Service struct {
    ID          string  `json:"id" gorm:"type:uuid;primaryKey"`
    Nom         string  `json:"nom"`
    Description string  `json:"description"`
    Prix        float64 `json:"prix"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
    if s.ID == "" {
        s.ID = uuid.New().String()
    }
    return nil
}

type Paiement struct {
    ID            string `json:"id" gorm:"type:uuid;primaryKey"`
    ReservationID string `json:"reservation_id" gorm:"type:uuid"`
    Montant       float64 `json:"montant"`
    ModePaiement  string `json:"mode_paiement" gorm:"type:varchar(50)"`
    Statut        string `json:"statut" gorm:"type:varchar(50);default:'en attente'"`
}

func (p *Paiement) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}