package main

import (
	"TestOzon/internal/config"
	"TestOzon/internal/repos/memory"
	"TestOzon/internal/repos/postgres"
	"TestOzon/internal/service"
	"context"
	"flag"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, opts))

	cfg := config.InitConfig()

	log.Info("Successfully initialized config")

	storageType := flag.String("storage", "memory", "type of storage to use: 'memory' or 'postgres'")
	flag.Parse()

	switch *storageType {
	case "postgres":
		db, err := postgres.InitPostgres(ctx, postgres.Config{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			DBName:   cfg.DB.DBName,
			SSLMode:  cfg.DB.SSLMode,
		})
		if err != nil {
			log.Error("Failed to initialize postgres storage", err.Error())
			os.Exit(1)
		}

		repos := postgres.NewReposPostgres(db)
		services := service.NewServiceP(repos)

	case "memory":
		db := memory.InitMemory()

		repos := memory.NewReposMemory(db)
		services := service.NewServiceM(repos)
	default:
		log.Error("Unsupported storage type", *storageType)
		log.Info("Use 'memory' or 'postgres'")
		os.Exit(1)
	}

}
