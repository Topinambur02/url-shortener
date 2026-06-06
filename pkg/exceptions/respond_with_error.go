package exceptions

import (
	"encoding/json"
	"net/http"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

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
