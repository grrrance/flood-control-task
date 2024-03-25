package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task/config"
	"task/pkg/logger"
)

type Server struct {
	cfg    *config.Config
	db     *redis.Client
	logger logger.Logger
}

const (
	maxHeaderBytes = 1 << 20
)

func NewServer(cfg *config.Config, db *redis.Client, logger logger.Logger) *Server {
	return &Server{cfg: cfg, db: db, logger: logger}
}

func (s *Server) Run() error {
	router := chi.NewRouter()
	err := s.MapHandlers(router)
	if err != nil {
		s.logger.Fatalf("Failed to configure the router: %s", err)
	}

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    s.cfg.Server.ReadTimeout,
		WriteTimeout:   s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		Handler:        router,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	go func() {
		<-sig
		s.logger.Info("Graceful shutdown")

		shutdownCtx, _ := context.WithTimeout(serverCtx, s.cfg.Server.CtxDefaultTimeout)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				s.logger.Fatalf("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Fatalf("Graceful shutdown has triggered with error: %s", err)
		}
		serverStopCtx()
	}()

	s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatalf("Error starting Server: %s", err)
		return err
	}

	<-serverCtx.Done()
	return nil
}
