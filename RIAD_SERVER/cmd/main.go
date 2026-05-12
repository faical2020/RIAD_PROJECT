package main

import (
	"log/slog"
	"net"
	"os"

	"RIAD_SERVER/internal/api"
	"RIAD_SERVER/internal/api/handlers"
	"RIAD_SERVER/internal/db"
	"RIAD_SERVER/internal/logic"
	"RIAD_SERVER/internal/sync"
	pb "github.com/anomalyco/riad_project/proto/sync"
	"google.golang.org/grpc"
)

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := logic.InitKeys(); err != nil {
		slogger.Error("failed to init RSA keys", "error", err)
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/riad?sslmode=disable"
	}

	if err := db.InitPostgres(databaseURL); err != nil {
		slogger.Error("db init failed", "error", err)
		os.Exit(1)
	}

	tokenStore := logic.NewRefreshTokenStore(db.GetDB())
	if err := tokenStore.AutoMigrate(); err != nil {
		slogger.Error("refresh token migration failed", "error", err)
		os.Exit(1)
	}
	handlers.SetTokenStore(tokenStore)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			slogger.Error("gRPC listen failed", "error", err)
			os.Exit(1)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterSyncServiceServer(grpcServer, sync.NewSyncServer())
		slogger.Info("gRPC sync server started", "port", 50051)
		if err := grpcServer.Serve(lis); err != nil {
			slogger.Error("gRPC serve failed", "error", err)
			os.Exit(1)
		}
	}()

	router := api.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slogger.Info("REST server started", "port", port)
	if err := router.Run(":" + port); err != nil {
		slogger.Error("server run failed", "error", err)
		os.Exit(1)
	}
}
