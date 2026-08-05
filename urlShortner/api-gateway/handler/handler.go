package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"system_design/urlShortner/api-gateway/models"
	"system_design/urlShortner/proto/url"
)

func CreateShortURL(client url.URLServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.UrlRequest
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.CreateShortURL(context.Background(), &url.CreateShortURLRequest{
			Url: event.URL,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetOriginalURL(client url.URLServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.URL.Query().Get("short_code")
		resp, err := client.GetOriginalURL(context.Background(), &url.GetOriginalURLRequest{
			ShortCode: shortCode,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
