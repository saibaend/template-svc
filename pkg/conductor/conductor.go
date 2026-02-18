package conductor

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Service interface {
	Name() string
	Run(chan<- error)
	Stop()
}

type Conductor struct {
	fifo []Service
	logg *slog.Logger
}

func New(logg *slog.Logger, fifo ...Service) Conductor {
	return Conductor{
		fifo: fifo,
		logg: logg,
	}
}

func (c Conductor) Run() chan error {
	errs := make(chan error, len(c.fifo))

	for _, svc := range c.fifo {
		go svc.Run(errs)
	}

	return errs
}

func (c Conductor) Shutdown(errs chan error) {
	interrupt := make(chan os.Signal, 1)

	go func() {
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	}()

	select {
	case <-interrupt:
		c.logg.LogAttrs(context.Background(), slog.LevelInfo, "conductor is shutting down by a signal")
	case err := <-errs:
		c.logg.LogAttrs(context.Background(), slog.LevelError, "conductor is stopping", slog.String("cause", err.Error()))
	}

	for _, svc := range c.fifo {
		svc.Stop()
		c.logg.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"service is stopped",
			slog.String("name", svc.Name()),
		)
	}
}
