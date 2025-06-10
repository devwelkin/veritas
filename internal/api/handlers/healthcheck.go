package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nouvadev/veritas/internal/config"
)

type API struct {
	App *config.AppConfig
}

func (cfg *API) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
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
