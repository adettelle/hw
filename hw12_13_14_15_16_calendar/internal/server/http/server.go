package internalhttp

import (
	"context"
	"net/http"

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

func NewServer(cfg *configs.Config, logg *zap.Logger, app Application) *Server {
	return &Server{cfg: cfg, logg: logg}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(s.cfg.Config, mux)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello!"))
}
