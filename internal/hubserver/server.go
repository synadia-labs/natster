package hubserver

import (
	log "log/slog"

	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/nats-io/nats.go"
)

type HubServer struct {
	nc      *nats.Conn
	library *medialibrary.MediaLibrary
}

func New(nc *nats.Conn, library *medialibrary.MediaLibrary) *HubServer {
	return &HubServer{
		nc:      nc,
		library: library,
	}
}

func (hub *HubServer) Start(uiPort int) error {
	err := hub.library.Ingest()
	if err != nil {
		return err
	}

	err = hub.startApiSubscriptions()
	if err != nil {
		return err
	}
	log.Info("Natster Media Hub Started")

	hub.startWebServer(uiPort)

	return nil
}

func (hub *HubServer) Stop() error {
	return nil
}
