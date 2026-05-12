package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_response "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set("X-Request-Id", requestID)
			w.Header().Set("X-Request-Id", requestID)

			ctx := WithRequestID(r.Context(), requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, _ := RequestIDFromContext(r.Context())
			if requestID == "" {
				requestID = r.Header.Get("X-Request-Id")
			}

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
			)

			ctx := context.WithValue(r.Context(), "logger", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())

			rw := &statusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			before := time.Now()

			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.statusCode),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
