package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

type URLHandler struct {
	s service.UrlService
}

func NewURLHandler(s service.UrlService) *URLHandler {
	return &URLHandler{s: s}
}

func (h *URLHandler) GetByShortUrl(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("short")

	if len(shortURL) != 10 {
		http.Error(w, "invalid short url length", http.StatusBadRequest)
		return
	}

	originalUrlDto, err := h.s.GetByShortUrl(r.Context(), shortURL)

	if err != nil {
		if err == exceptions.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(originalUrlDto)
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUrlDto dto.CreateUrlDto

	if err := json.NewDecoder(r.Body).Decode(&createUrlDto); err != nil || createUrlDto.OriginalUrl == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	shortUrlDto, err := h.s.Create(r.Context(), &createUrlDto)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortUrlDto)
}