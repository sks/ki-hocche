package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/sks/kihocche/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Port                    string
	ShutdownTimeoutDuration time.Duration
}

func (c Config) Start(ctx context.Context, handler http.Handler) error {
	errGroup, ectx := errgroup.WithContext(ctx)
	server := http.Server{
		Addr: net.JoinHostPort("", c.Port),
		BaseContext: func(listener net.Listener) context.Context {
			return ectx
		},
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errGroup.Go(func() error {
		logger.GetLogger(ctx).Info("server started", "port", c.Port)
		return server.ListenAndServe()
	})
	errGroup.Go(func() error {
		<-ctx.Done()
		shutDownCtx, cancel := context.WithTimeout(context.Background(), c.ShutdownTimeoutDuration)
		defer cancel()
		return server.Shutdown(shutDownCtx)
	})
	return errGroup.Wait()
}
