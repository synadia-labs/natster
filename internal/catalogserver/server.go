package catalogserver

import (
	log "log/slog"

	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/nats-io/nats.go"
)

type CatalogServer struct {
	nc      *nats.Conn
	library *medialibrary.MediaLibrary
}

func New(nc *nats.Conn, library *medialibrary.MediaLibrary) *CatalogServer {
	return &CatalogServer{
		nc:      nc,
		library: library,
	}
}

func (srv *CatalogServer) Start(uiPort int) error {
	err := srv.startApiSubscriptions()
	if err != nil {
		return err
	}
	log.Info("Natster Media Catalog Server Started")

	srv.startWebServer(uiPort)

	return nil
}

func (srv *CatalogServer) Stop() error {
	return nil
}
