package middleware

import (
	"net/http"
	"time"

	"github.com/topinambur02/url-shortener/pkg/logging"
	"github.com/topinambur02/url-shortener/pkg/utils"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.size += int64(n)
	return n, err
}

func HTTPLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		requestID := r.Header.Get("X-Request-ID")

		if requestID != "" {
			logger = logger.GetLoggerWithField("request_id", requestID)
		}

		logger = logger.GetLoggerWithField("method", r.Method)
		logger = logger.GetLoggerWithField("path", r.URL.Path)
		logger = logger.GetLoggerWithField("client_ip", utils.GetClientIP(r))
		logger = logger.GetLoggerWithField("user_agent", r.UserAgent())

		start := time.Now()

		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		latency := time.Since(start)

		logger = logger.GetLoggerWithField("status", lrw.statusCode)
		logger = logger.GetLoggerWithField("latency_ms", latency.Milliseconds())
		logger = logger.GetLoggerWithField("size", lrw.size)

		switch {
		case lrw.statusCode >= 500:
			logger.Error("Request completed with server error")
		case lrw.statusCode >= 400:
			logger.Warn("Request completed with client error")
		default:
			logger.Info("Request processed")
		}
	})
}
