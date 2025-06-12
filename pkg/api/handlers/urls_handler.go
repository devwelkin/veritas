package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
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

type RedirectEvent struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
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
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Short code is required")
		return
	}

	// 1. Try to get from cache first
	originalURL, err := h.App.Cache.Get(r.Context(), shortCode).Result()
	if err == nil {
		h.App.Logger.Info("cache hit", "short_code", shortCode)
		// Redirect and publish event
		h.publishRedirectEvent(shortCode, originalURL, r)
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	if !errors.Is(err, redis.Nil) {
		h.App.Logger.Error("redis error", "err", err)
	} else {
		h.App.Logger.Info("cache miss", "short_code", shortCode)
	}

	// 2. If not in cache, get from DB
	originalURL, err = h.App.Querier.GetURLByShortCode(r.Context(), shortCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "URL not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get URL")
		}
		h.App.Logger.Error("db error", "err", err)
		return
	}

	// 3. Store in cache for future requests
	if err := h.App.Cache.Set(r.Context(), shortCode, originalURL, 1*time.Hour).Err(); err != nil {
		h.App.Logger.Error("failed to set cache", "err", err)
	}

	// Redirect and publish event
	h.publishRedirectEvent(shortCode, originalURL, r)
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *URLHandler) publishRedirectEvent(shortCode, originalURL string, r *http.Request) {
	event := RedirectEvent{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserAgent:   r.UserAgent(),
		IPAddress:   r.RemoteAddr,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		h.App.Logger.Error("failed to marshal redirect event", "err", err)
		return
	}

	subject := "veritas.redirect.success"
	if err := h.App.NATS.Publish(subject, eventJSON); err != nil {
		h.App.Logger.Error("failed to publish nats event", "err", err)
	} else {
		h.App.Logger.Info("published nats event", "subject", subject)
	}
}
