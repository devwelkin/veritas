package utils

import (
	"net/http"
	"net/url"
	"time"
)

func ValidateURL(urlToTest string) bool {
	// seviye 1: format kontrolü
	_, err := url.ParseRequestURI(urlToTest)
	if err != nil {
		return false // formatı bile bozuk, hiç uğraşma.
	}

	// ulaşılabilirlik kontrolü
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(urlToTest)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 400'den küçük bir status code (2xx, 3xx) yeterli.
	return resp.StatusCode < 400
}
