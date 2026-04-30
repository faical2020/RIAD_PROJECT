package db

import (
    "fmt"
    "log"
    "RIAD_SERVER/internal/logic"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres(databaseURL string) error {
    var err error
    DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("échec connexion PostgreSQL: %w", err)
    }

    log.Println("Migration des modèles...")
    err = DB.AutoMigrate(
        &logic.User{},
        &logic.Chambre{},
        &logic.Reservation{},
        &logic.Tache{},
        &logic.Service{},
        &logic.Paiement{},
    )
    if err != nil {
        return fmt.Errorf("échec migration: %w", err)
    }

    log.Println("Base PostgreSQL initialisée")
    return nil
}

func GetDB() *gorm.DB {
    return DB
}