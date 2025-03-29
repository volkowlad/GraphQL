package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func InitPostgres(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	err := cfg.Validate()
	if err != nil {
		slog.Error("Error validating postgres configuration", err.Error())
		return nil, err
	}

	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		slog.Error("error connecting to postgres", err.Error())
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		slog.Error("error to ping postgres", err.Error())
		return nil, err
	}

	return db, nil
}
