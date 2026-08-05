package handler

import (
	"context"
	"log"
	"system_design/urlShortner/internal/cache"
	"system_design/urlShortner/internal/repository"
	"system_design/urlShortner/internal/utils"
	"system_design/urlShortner/proto/url"
)

type URLHandler struct {
	url.UnimplementedURLServiceServer
	repo  *repository.URLRepository
	cache *cache.RedisCache
}

func NewURLHandler(repo *repository.URLRepository, cache *cache.RedisCache) *URLHandler {
	return &URLHandler{repo: repo, cache: cache}
}
func (h *URLHandler) CreateShortURL(ctx context.Context, req *url.CreateShortURLRequest) (*url.CreateShortURLResponse, error) {
	shortCode := utils.GenerateShortCode()

	h.repo.CreateShortURL(req.Url, shortCode)
	h.cache.Set(shortCode, req.Url)
	return &url.CreateShortURLResponse{
		ShortCode: shortCode,
		ShortUrl:  "http://localhost:8080/" + shortCode,
	}, nil

}
func (h *URLHandler) GetOriginalURL(ctx context.Context, req *url.GetOriginalURLRequest) (*url.GetOriginalURLResponse, error) {
	var originalURL string

	data, err := h.cache.Get(req.ShortCode)
	if err == nil {
		log.Println("Redis HIT")
		originalURL = data
	} else {
		log.Println("Redis MISS")

		originalURL, err = h.repo.GetOriginalURL(req.ShortCode)
		if err != nil {
			return nil, err
		}

		if err := h.cache.Set(req.ShortCode, originalURL); err != nil {
			log.Println("Failed to cache:", err)
		}
	}

	return &url.GetOriginalURLResponse{
		OriginalUrl: originalURL,
	}, nil
}
