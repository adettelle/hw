package logger

import (
	"net"
	"net/http"
	"strings"
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
	// если статус ещё не установлен, считаем его 200
	if r.responseData.status == 0 {
		r.responseData.status = http.StatusOK
	}
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

		if responseData.status == 0 {
			responseData.status = http.StatusOK
		}
		latencyMs := float64(duration.Nanoseconds())

		logger.Info("http_request",
			zap.String("cient_ip", clientIP(r)),
			zap.Time("time", time.Now()),
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("proto", r.Proto),
			zap.Int("status", responseData.status),
			zap.Float64("latency_ms", latencyMs),
			zap.Duration("duration", duration),
			zap.Int("response_size", responseData.size),
			zap.String("user_agent", r.UserAgent()),
		)
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				return p
			}
		}
	}

	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
