package api

import (
	"net/http"

	"github.com/nouvadev/veritas/internal/api/handlers"
	"github.com/nouvadev/veritas/internal/config"
)

// Routes sets up the routes for the application.
func Routes(app *config.AppConfig) http.Handler {
	mux := http.NewServeMux()

	h := handlers.NewHealthcheckHandler(app)
	u := handlers.NewURLHandler(app)

	mux.HandleFunc("GET /healthcheck", h.HealthcheckHandler)
	mux.HandleFunc("POST /create", u.CreateShortURL)

	return mux
}
