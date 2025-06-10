package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nouvadev/veritas/internal/config"
)

// URLHandler handles all URL-related HTTP requests
type URLHandler struct {
	App *config.AppConfig
}

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

func NewURLHandler(app *config.AppConfig) *URLHandler {
	return &URLHandler{App: app}
}

func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}
