package conductor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	srvHttp "github.com/infoland-kz/step2travel_manager_api/pkg/http"
)

const (
	port = ":9999"
	host = "http://localhost:9999"

	keyWord = "pong"
)

type pingPongHandler struct{}

func (h pingPongHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, keyWord)
}

func TestConductor_Run(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "success",
			f: func(t *testing.T) {
				logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

				handler := pingPongHandler{}
				srv := srvHttp.New(port, logger, srvHttp.WithRouter(handler))

				go func() {
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				cdc := New(logger, srv)
				cdc.Shutdown(cdc.Run())
			},
		},
		{
			name: "server already started",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				handler := pingPongHandler{}
				srv := srvHttp.New(port, logger, srvHttp.WithRouter(handler))

				srv2 := srvHttp.New(port, logger, srvHttp.WithRouter(handler))

				go func() {
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				}()

				cdc := New(logger, srv, srv2)
				cdc.Shutdown(cdc.Run())

				if !bytes.Contains(buff.Bytes(), []byte("listen tcp :9999: bind: address already in use")) {
					t.Errorf("unexpected error message: %s", buff.String())
				}
			},
		},
		{
			name: "graceful shutdown",
			f: func(t *testing.T) {
				buff := &bytes.Buffer{}

				logger := slog.New(slog.NewJSONHandler(buff, nil))

				srv := srvHttp.New(port, logger, srvHttp.WithRouter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					_, _ = fmt.Fprintf(w, keyWord)
				})))

				first := make(chan struct{})
				second := make(chan struct{})
				third := make(chan struct{})

				go func() {
					<-first
					go func() {
						second <- struct{}{}
					}()
					res, err := http.Get(host)
					if err != nil {
						t.Errorf("failed to send first request: %v", err)
					}

					data, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("failed to parse reposne: %v", err)
					}

					if string(data) != keyWord {
						t.Errorf("unexpected key word: got %v", string(data))
					}
				}()

				go func() {
					<-third
					_, err := http.Get(host)
					if err == nil {
						t.Errorf("expected error as second request response, got nil")
						return
					}

					if !errors.Is(err, syscall.ECONNREFUSED) {
						t.Errorf("unexpected error, got %v, want %v", err, syscall.ECONNREFUSED)
					}
				}()

				go func() {
					<-second
					time.Sleep(1 * time.Second)
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					third <- struct{}{}
				}()

				go func() {
					first <- struct{}{}
				}()
				cdc := New(logger, srv)
				cdc.Shutdown(cdc.Run())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.f)
	}
}
