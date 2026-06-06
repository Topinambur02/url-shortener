package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

var validate = validator.New()

type URLHandler struct {
	s       service.URLService
	address string
}

func NewURLHandler(s service.URLService, address string) *URLHandler {
	return &URLHandler{s: s, address: address}
}

// GetByShortUrl godoc
// @Summary      Получить оригинальный URL
// @Description  Возвращает оригинальный URL по его 10-символьному короткому коду.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        short   path      string  true  "Короткий код (10 символов)" minlength(10) maxlength(10)
// @Success      200     {object}  dto.OriginalURLDto
// @Failure      400     {object}  exceptions.ErrBadRequestResponseDoc
// @Failure      404     {object}  exceptions.ErrNotFoundResponseDoc
// @Failure      500     {object}  exceptions.ErrInternalResponseDoc
// @Router       /{short} [get]
func (h *URLHandler) GetByShortURL(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	shortURL := r.PathValue("short")

	if len(shortURL) != 10 {
		logger.Infof("Error: Invalid shortURL length (%d characters): %s", len(shortURL), shortURL)
		exceptions.RespondWithError(w, http.StatusBadRequest, "invalid short url length")
		return
	}

	originalUrlDto, err := h.s.GetByShortURL(r.Context(), shortURL)

	if err != nil {
		if err == exceptions.ErrNotFound {
			logger.Infof("URL not found in database: %s", shortURL)
			exceptions.RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		
		logger.Infof("Internal error while searching %s: %v", shortURL, err)
		exceptions.RespondWithError(w, http.StatusInternalServerError, "internal error")
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
// @Param        request body      dto.CreateURLDto  true  "Данные для создания ссылки"
// @Success      201     {object}  dto.ShortURLDto
// @Failure      400     {object}  exceptions.ErrBadRequestResponseDoc
// @Failure		 409	 {object}  exceptions.ErrConflictResponseDoc
// @Failure      500     {object}  exceptions.ErrInternalResponseDoc
// @Router       / [post]
func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	var createUrlDto dto.CreateURLDto

	if err := json.NewDecoder(r.Body).Decode(&createUrlDto); err != nil {
		logger.Infof("Error decoding request body: %v", err)
		exceptions.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := validate.Struct(createUrlDto); err != nil {
		exceptions.RespondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	shortUrlDto, err := h.s.Create(r.Context(), &createUrlDto, h.address)

	if err != nil {
        if errors.Is(err, exceptions.ErrConflict) {
            logger.Warnf("Conflict: URL already shortened: %s", createUrlDto.OriginalURL)
            exceptions.RespondWithError(w, http.StatusConflict, "url already exists")
            return
        }

        logger.Errorf("Internal error while creating short link for %s: %v", createUrlDto.OriginalURL, err)
        exceptions.RespondWithError(w, http.StatusInternalServerError, "internal error")
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(shortUrlDto); err != nil {
		logger.Infof("JSON response encoding error: %v", err)
		return
	}
}
