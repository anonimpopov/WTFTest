package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	*http.Server
}

func New(addr string, handler http.Handler) (*server, <-chan os.Signal) {
	srv := &server{
		Server: &http.Server{
			Addr:        addr,
			Handler:     handler,
			ReadTimeout: 3 * time.Second,
			IdleTimeout: 65 * time.Second,
		},
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	return srv, shutdown
}

func (srv *server) Shutdown(ctx context.Context) {
	srv.SetKeepAlivesEnabled(false)
	if err := srv.Server.Shutdown(ctx); err != nil {
		logrus.Errorf("error while shutdowning server: %v", err)
	}

	_ = srv.Close()
}
