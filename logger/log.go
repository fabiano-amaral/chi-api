package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var Logger *zap.Logger

type Config struct {
	LogReferer   bool
	LogUserAgent bool
}

func Init() {
	Logger, _ = zap.NewProduction()
}

func RequestMiddleware(logger *zap.Logger, c *Config) func(next http.Handler) http.Handler {
	if logger == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				reqLogger := logger.With(
					zap.String("path", r.URL.Path),
					zap.String("requestID", middleware.GetReqID(r.Context())),
					zap.Duration("latency", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
				if c.LogReferer {
					referer := ww.Header().Get("Referer")
					if referer == "" {
						referer = r.Header.Get("Referer")
					}
					if referer != "" {
						reqLogger = reqLogger.With(zap.String("ref", referer))
					}
				}
				if c.LogUserAgent {
					userAgent := ww.Header().Get("User-Agent")
					if userAgent == "" {
						userAgent = r.Header.Get("User-Agent")
					}
					if userAgent != "" {
						reqLogger = reqLogger.With(zap.String("user-agent", userAgent))
					}
				}
				reqLogger.Info("Served")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
