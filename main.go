package main

import (
	"net/http"
	"time"

	log "github.com/fabiano-amaral/chi-api/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
  r := chi.NewRouter()
  logger, _ := zap.NewProduction()
  r.Use(log.RequestMiddleware(logger, &log.Config{
    LogReferer: true,
    LogUserAgent: true,
  }))
  r.Use(middleware.RequestID)
  r.Use(middleware.Recoverer)
  r.Use(middleware.StripSlashes)
  r.Use(middleware.Timeout(60 * time.Second))
  r.Get("/status", func(w http.ResponseWriter, _ *http.Request)  {
    w.Write([]byte("Olá mundo!!"))
  })

  http.ListenAndServe(":3000", r)
}
