package utils

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// ValidateURL checks if a URL has a valid format (scheme and host).
// It does not check for reachability.
func ValidateURL(urlToTest string) bool {
	// Simple format check using Go's standard library.
	// This is sufficient for a URL shortener as we don't need to
	// guarantee reachability, only that the URL is well-formed.
	parsedURL, err := url.ParseRequestURI(urlToTest)
	if err != nil {
		return false
	}

	// Ensure the URL has a scheme (http, https) and a host.
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// CheckURLReachability performs a HEAD request to see if a URL is reachable.
func CheckURLReachability(urlToCheck string, logger *slog.Logger) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(urlToCheck)
	if err != nil {
		logger.Warn("Reachability check failed", "url", urlToCheck, "error", err)
		return false
	}
	defer resp.Body.Close()

	// Any status code less than 400 (e.g., 2xx or 3xx) is considered successful.
	isReachable := resp.StatusCode < 400
	if !isReachable {
		logger.Warn("Reachability check returned non-success status", "url", urlToCheck, "status_code", resp.StatusCode)
	}
	return isReachable
}
