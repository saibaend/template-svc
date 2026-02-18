package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Option func(s *Server)

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readHeaderTimeout = timeout
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func WithRouter(handler http.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}

func WithTLSConfig(config *tls.Config) Option {
	return func(s *Server) {
		s.tlsConfig = config
	}
}
