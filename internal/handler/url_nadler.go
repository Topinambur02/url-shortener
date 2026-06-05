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

// GetByShortUrl godoc
// @Summary      Получить оригинальный URL
// @Description  Возвращает оригинальный URL по его 10-символьному короткому коду.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        short   path      string  true  "Короткий код (10 символов)" minlength(10) maxlength(10)
// @Success      200     {object}  dto.OriginalUrlDto
// @Failure      400     {string}  string  "invalid short url length"
// @Failure      404     {string}  string  "not found"
// @Failure      500     {string}  string  "internal error"
// @Router       /{short} [get]
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
	if err := json.NewEncoder(w).Encode(originalUrlDto); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// Create godoc
// @Summary      Создать короткую ссылку
// @Description  Принимает оригинальный URL в теле запроса и генерирует для него короткий код.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        request body      dto.CreateUrlDto  true  "Данные для создания ссылки"
// @Success      200     {object}  dto.ShortUrlDto
// @Failure      400     {string}  string  "invalid request"
// @Failure      500     {string}  string  "internal error"
// @Router       / [post]
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
	if err := json.NewEncoder(w).Encode(shortUrlDto); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
