package api

import (
	"database/sql"
	"log/slog"
)

// application struct holds the dependencies for our HTTP handlers, helpers, and middleware.
type Application struct {
	Logger *slog.Logger
}

type App struct {
	log *slog.Logger
	db  *sql.DB
}
