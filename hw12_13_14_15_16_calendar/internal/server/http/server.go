package internalhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/configs"
	"go.uber.org/zap"
)

type Server struct { // TODO
	cfg  *configs.Config
	logg *zap.Logger
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(cfg *configs.Config, logg *zap.Logger, _ Application) *Server {
	return &Server{cfg: cfg, logg: logg}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         s.cfg.Config,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	mux.HandleFunc(`/`, mainPage)

	// err := http.ListenAndServe(s.cfg.Config, mux)
	// if err != nil {
	// 	panic(err)
	// }
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	// TODO
	return nil
}

func mainPage(res http.ResponseWriter, _ *http.Request) {
	res.Write([]byte("Hello!"))
}
