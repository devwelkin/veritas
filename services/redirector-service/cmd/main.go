package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nouvadev/veritas/pkg/api/handlers"
	"github.com/nouvadev/veritas/pkg/cache"
	"github.com/nouvadev/veritas/pkg/config"
	"github.com/nouvadev/veritas/pkg/database"
	sqlc "github.com/nouvadev/veritas/pkg/database/sqlc"
	"github.com/nouvadev/veritas/pkg/nats"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file, continuing without it...")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	dbpool, err := database.ConnectDB(dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	logger.Info("database connection pool established")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		logger.Error("REDIS_URL environment variable is not set")
		os.Exit(1)
	}

	redisClient, err := cache.ConnectRedis(redisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "err", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	logger.Info("redis connection established")

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		logger.Error("NATS_URL environment variable is not set")
		os.Exit(1)
	}

	natsConn, err := nats.ConnectNATS(natsURL)
	if err != nil {
		logger.Error("failed to connect to nats", "err", err)
		os.Exit(1)
	}
	defer natsConn.Close()

	queries := sqlc.New(dbpool)

	app := &config.AppConfig{
		Logger:  logger,
		DB:      dbpool,
		Querier: queries,
		Cache:   redisClient,
		NATS:    natsConn,
	}

	PORT := os.Getenv("REDIRECTOR_PORT")
	if PORT == "" {
		PORT = "8082"
	}
	logger.Info("starting server", "addr", PORT)

	mux := http.NewServeMux()

	h := handlers.NewHealthcheckHandler(app)
	u := handlers.NewURLHandler(app)

	mux.HandleFunc("GET /healthcheck", h.HealthcheckHandler)
	mux.HandleFunc("GET /{short_code}", u.RedirectToOriginalURL)

	err = http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
