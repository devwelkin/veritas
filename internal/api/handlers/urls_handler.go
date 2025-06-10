package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nouvadev/veritas/internal/config"
	database "github.com/nouvadev/veritas/internal/database/sqlc"
	"github.com/nouvadev/veritas/internal/utils"
)

// URLHandler handles all URL-related HTTP requests
type URLHandler struct {
	App *config.AppConfig
}

type URLRequest struct {
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

func NewURLHandler(app *config.AppConfig) *URLHandler {
	return &URLHandler{App: app}
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// get original url from request
	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		h.App.Logger.Error("Invalid request body", "error", err)
		return
	}

	if !utils.ValidateURL(req.OriginalURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		h.App.Logger.Error("Invalid URL", "url", req.OriginalURL)
		return
	}

	insertedID, err := h.App.Querier.CreateURL(r.Context(), req.OriginalURL)
	if err != nil {
		http.Error(w, "Failed to create URL", http.StatusInternalServerError)
		h.App.Logger.Error("Failed to create URL", "error", err)
		return
	}

	shortCode := utils.ToBase62(uint64(insertedID))

	err = h.App.Querier.UpdateShortCode(r.Context(), database.UpdateShortCodeParams{
		ShortCode: shortCode,
		ID:        insertedID,
	})
	if err != nil {
		http.Error(w, "Failed to update short code", http.StatusInternalServerError)
		h.App.Logger.Error("Failed to update short code", "error", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, URLResponse{ShortURL: shortCode})
}
