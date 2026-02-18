package http

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout       = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultReadHeaderTimeout = 10 * time.Second
	defaultShutdownTimeout   = 30 * time.Second
)

type Server struct {
	addr              string
	handler           http.Handler
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	shutdownTimeout   time.Duration

	logg      *slog.Logger
	srv       *http.Server
	tlsConfig *tls.Config
}

func New(addr string, logg *slog.Logger, opts ...Option) Server {
	//nolint:exhaustruct
	s := Server{
		addr: addr,
		logg: logg,
	}

	for _, opt := range opts {
		opt(&s)
	}

	//nolint:exhaustruct
	s.srv = &http.Server{
		Addr:              addr,
		Handler:           s.handler,
		ReadHeaderTimeout: s.readHeaderTimeout,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
		TLSConfig:         s.tlsConfig,
	}

	if s.writeTimeout <= 0 {
		s.writeTimeout = defaultWriteTimeout
	}

	if s.readTimeout <= 0 {
		s.readTimeout = defaultReadTimeout
	}

	if s.readHeaderTimeout <= 0 {
		s.readHeaderTimeout = defaultReadHeaderTimeout
	}

	if s.shutdownTimeout <= 0 {
		s.shutdownTimeout = defaultShutdownTimeout
	}

	return s
}

func (s Server) Run(errs chan<- error) {
	s.logg.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"running server",
		slog.String("server host", s.addr),
	)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- err
	}

	errs <- nil
}

func (s Server) Stop() {
	sCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(sCtx); err != nil {
		s.logg.LogAttrs(
			context.Background(),
			slog.LevelError,
			"server shutting down error",
			slog.String("server host", s.addr),
			slog.Any("error", err),
		)
	}
}

func (s Server) Name() string {
	return s.addr
}
