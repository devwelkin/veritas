package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/nouvadev/veritas/pkg/config"
	database "github.com/nouvadev/veritas/pkg/database/sqlc"
	"github.com/nouvadev/veritas/pkg/utils"
	"github.com/redis/go-redis/v9"
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		h.App.Logger.Error("Invalid request body", "error", err)
		return
	}

	if !utils.ValidateURL(req.OriginalURL) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL")
		h.App.Logger.Error("Invalid URL", "url", req.OriginalURL)
		return
	}

	insertedID, err := h.App.Querier.CreateURL(r.Context(), req.OriginalURL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create URL")
		h.App.Logger.Error("Failed to create URL", "error", err)
		return
	}

	shortCode := utils.ToBase62(uint64(insertedID))

	err = h.App.Querier.UpdateShortCode(r.Context(), database.UpdateShortCodeParams{
		ShortCode: shortCode,
		ID:        insertedID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update short code")
		h.App.Logger.Error("Failed to update short code", "error", err)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), shortCode)

	// Respond to the user immediately
	utils.RespondWithJSON(w, http.StatusCreated, URLResponse{ShortURL: shortURL})

	// Perform reachability check in the background
	go func() {
		isReachable := utils.CheckURLReachability(req.OriginalURL, h.App.Logger)
		if !isReachable {
			h.App.Logger.Info("URL is not reachable, deleting", "id", insertedID)
			err := h.App.Querier.DeleteURL(r.Context(), insertedID)
			if err != nil {
				h.App.Logger.Error("Failed to delete unreachable URL", "id", insertedID, "error", err)
			}
		}
	}()
}

func (h *URLHandler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]

	originalURL, err := h.App.Cache.Get(r.Context(), shortCode).Result()
	if err == nil {
		h.App.Logger.Info("URL found in cache", "short_code", shortCode, "original_url", originalURL)
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
		return
	}

	if !errors.Is(err, redis.Nil) {
		h.App.Logger.Error("redis error", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.App.Logger.Info("URL not found in cache", "short_code", shortCode)

	originalURL, err = h.App.Querier.GetURLByShortCode(r.Context(), shortCode)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "URL not found")
		h.App.Logger.Error("URL not found", "short_code", shortCode)
		return
	}

	err = h.App.Cache.Set(r.Context(), shortCode, originalURL, 0).Err()
	if err != nil {
		h.App.Logger.Error("failed to set URL in cache", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
