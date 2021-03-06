package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(address string, logging bool, logger logging.Logger,
	buildInfo models.BuildInformation,
	openvpnLooper openvpn.Looper, unboundLooper dns.Looper,
	updaterLooper updater.Looper, publicIPLooper publicip.Looper) Server {
	serverLogger := logger.WithPrefix("http server: ")
	handler := newHandler(serverLogger, logging, buildInfo,
		openvpnLooper, unboundLooper, updaterLooper, publicIPLooper)
	return &server{
		address: address,
		logger:  serverLogger,
		handler: handler,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("context canceled: shutting down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()
	s.logger.Info("listening on %s", s.address)
	err := server.ListenAndServe()
	if err != nil && ctx.Err() != context.Canceled {
		s.logger.Error(err)
	}
	s.logger.Warn("shut down")
}
