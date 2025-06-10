package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nouvadev/veritas/internal/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &api.Application{
		Logger: logger,
	}

	addr := ":8080"
	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, app.Routes())
	if err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
