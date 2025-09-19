package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData is a structure used to store information about an HTTP response.
type responseData struct {
	status int // HTTP status code of the response
	size   int // Size of the response body in bytes
}

// добавляем реализацию http.ResponseWriter.
type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// WithLogging wraps an http.HandlerFunc to add logging functionality.
// It logs information about each HTTP request, including the URI, method, response status code,
// response size, and the duration of the request.
func WithLogging(h http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{status: 0, size: 0}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		logger.Info("Request data:", zap.String("uri", r.RequestURI),
			zap.String("method", r.Method), zap.Int("ststus", responseData.status),
			zap.Duration("duration", duration), zap.Int("size", responseData.size))
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
