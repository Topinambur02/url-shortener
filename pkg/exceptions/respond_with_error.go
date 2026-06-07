package exceptions

import (
	"encoding/json"
	"net/http"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

type ErrBadRequestResponse struct {
    dto.ErrorDto
}

type ErrBadRequestResponseDoc struct {
    StatusCode int    `json:"status_code" example:"400"`
    Message    string `json:"message"     example:"invalid request"`
}

type ErrNotFoundResponseDoc struct {
	StatusCode int    `json:"status_code" example:"404"`
    Message    string `json:"message"     example:"url not found"`
}

type ErrConflictResponseDoc struct {
    StatusCode int    `json:"status_code" example:"409"`
    Message    string `json:"message"     example:"url already exists"`
}

type ErrInternalResponseDoc struct {
    StatusCode int    `json:"status_code" example:"500"`
    Message    string `json:"message"     example:"internal error"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	logger := logging.GetLogger()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    errDto := dto.ErrorDto{
        StatusCode: statusCode,
        Message:    message,
    }
    
	if err := json.NewEncoder(w).Encode(errDto); err != nil {
		logger.Infof("Error encoding JSON response for: %s", err)
		return
	}
}
