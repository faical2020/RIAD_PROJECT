package main

import (
    "log"
    "os"
    "RIAD_SERVER/internal/api"
    "RIAD_SERVER/internal/db"
)

func main() {
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        databaseURL = "postgres://postgres:postgres@localhost:5432/riad?sslmode=disable"
    }

    if err := db.InitPostgres(databaseURL); err != nil {
        log.Fatal("Erreur DB:", err)
    }

    router := api.SetupRouter()
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Serveur démarré sur :%s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Erreur serveur:", err)
    }
}