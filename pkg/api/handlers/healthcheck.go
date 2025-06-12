package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nouvadev/veritas/pkg/config"
)

type HealthcheckHandler struct {
	App *config.AppConfig
}

func NewHealthcheckHandler(app *config.AppConfig) *HealthcheckHandler {
	return &HealthcheckHandler{App: app}
}

func (h *HealthcheckHandler) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status": "ok",
	}

	js, err := json.Marshal(health)
	if err != nil {
		http.Error(w, "error marshalling healthcheck", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
