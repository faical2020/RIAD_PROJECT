package main

import (
	"log"
	"net"
	"os"

	"RIAD_SERVER/internal/api"
	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/sync"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/riad?sslmode=disable"
	}

	if err := db.InitPostgres(databaseURL); err != nil {
		log.Fatal("Erreur DB:", err)
	}

	// Start gRPC Server for Synchronization
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterSyncServiceServer(grpcServer, sync.NewSyncServer())
		log.Printf("gRPC Sync Server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	router := api.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Serveur REST démarré sur :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}
