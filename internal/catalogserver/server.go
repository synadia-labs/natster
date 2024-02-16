package catalogserver

import (
	log "log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/globalservice"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	heartbeatIntervalSeconds = 30
)

type CatalogServer struct {
	nc                  *nats.Conn
	globalServiceClient *globalservice.Client
	library             *medialibrary.MediaLibrary
	nctx                *models.NatsterContext
	hbQuit              chan bool
}

func New(ctx *models.NatsterContext, nc *nats.Conn, library *medialibrary.MediaLibrary) *CatalogServer {
	return &CatalogServer{
		nc:                  nc,
		hbQuit:              make(chan bool),
		nctx:                ctx,
		globalServiceClient: globalservice.NewClient(nc),
		library:             library,
	}
}

func (srv *CatalogServer) Start(uiPort int) error {
	err := srv.startApiSubscriptions()
	if err != nil {
		return err
	}
	log.Info("Natster Media Catalog Server Started")

	srv.startHeartbeatEmitter()
	srv.startWebServer(uiPort)

	return nil
}

func (srv *CatalogServer) startHeartbeatEmitter() {
	ticker := time.NewTicker(heartbeatIntervalSeconds * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				_ = srv.globalServiceClient.PublishHeartbeat(srv.nctx, srv.library.Name)
			case <-srv.hbQuit:
				ticker.Stop()
				close(srv.hbQuit)
				return
			}
		}
	}()
}

func (srv *CatalogServer) Stop() error {
	srv.hbQuit <- true
	return nil
}
