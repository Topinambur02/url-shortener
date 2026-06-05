package handler

import (
	"encoding/json"
	"net/http"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
	"github.com/topinambur02/url-shortener/pkg/logging"
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
	logger := logging.GetLogger()
	shortURL := r.PathValue("short")

	if len(shortURL) != 10 {
		logger.Infof("Error: Invalid shortURL length (%d characters): %s", len(shortURL), shortURL)
		http.Error(w, "invalid short url length", http.StatusBadRequest)
		return
	}

	originalUrlDto, err := h.s.GetByShortUrl(r.Context(), shortURL)

	if err != nil {
		if err == exceptions.ErrNotFound {
			logger.Infof("URL not found in database: %s", shortURL)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		logger.Infof("Internal error while searching %s: %v", shortURL, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(originalUrlDto); err != nil {
		logger.Infof("Error encoding JSON response for %s: %v", shortURL, err)
		return
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
	logger := logging.GetLogger()
    var createUrlDto dto.CreateUrlDto

    if err := json.NewDecoder(r.Body).Decode(&createUrlDto); err != nil {
        logger.Infof("Error decoding request body: %v", err)
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if createUrlDto.OriginalUrl == "" {
        logger.Infoln("error: empty field OriginalUrl")
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    shortUrlDto, err := h.s.Create(r.Context(), &createUrlDto)

    if err != nil {
        logger.Infof("Internal error while creating short link for %s: %v", createUrlDto.OriginalUrl, err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(shortUrlDto); err != nil {
        logger.Infof("JSON response encoding error: %v", err)
        return
    }
}
