package config

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AppConfig struct holds the dependencies for our HTTP handlers, helpers, and middleware.
type AppConfig struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
}
