package catalogserver

import (
	log "log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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
	allowAll            bool
	catalogWatcher      *fsnotify.Watcher
}

func New(ctx *models.NatsterContext, nc *nats.Conn, library *medialibrary.MediaLibrary, allowAll bool) *CatalogServer {
	return &CatalogServer{
		nc:                  nc,
		hbQuit:              make(chan bool),
		nctx:                ctx,
		globalServiceClient: globalservice.NewClient(nc),
		library:             library,
		allowAll:            allowAll,
		catalogWatcher:      nil,
	}
}

func (srv *CatalogServer) Start() error {
	if srv.allowAll {
		log.Warn("WARNING - This server will not enforce any security checks on contents queries or downloads.")
		log.Warn("If your catalogs are exported, anyone with your 56-character account ID will be able to acccess your media.")
		log.Warn("If this is not what you wanted, please stop this server immediately.")
	}
	err := srv.startApiSubscriptions()
	if err != nil {
		return err
	}
	log.Info("Natster Media Catalog Server Started")
	log.Info("Local (private) services are available on 'natster.local.>'")

	srv.startHeartbeatEmitter()
	srv.startCatalogMonitor()

	return nil
}

func (srv *CatalogServer) startHeartbeatEmitter() {
	ticker := time.NewTicker(heartbeatIntervalSeconds * time.Second)
	// publish one immediately
	_ = srv.globalServiceClient.PublishHeartbeat(srv.nctx, srv.library.Name)

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

func (srv *CatalogServer) startCatalogMonitor() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Warn("Failed to create a file system watcher for the catalog. Data will not live update", err)
		return
	}
	// Start listening for events.
	go srv.watchLoop(watcher)

	_ = watcher.Add(srv.library.RootDir)

	for _, entry := range srv.library.Entries {
		theDir := filepath.Dir(entry.Path)
		_ = watcher.Add(theDir)
	}

	if err != nil {
		log.Warn("Failed to watch catalog root directory", err)
		return
	}

}

func (srv *CatalogServer) Stop() error {
	srv.hbQuit <- true
	if srv.catalogWatcher != nil {
		_ = srv.catalogWatcher.Close()
	}

	return nil
}

func (srv *CatalogServer) watchLoop(w *fsnotify.Watcher) {
	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Error("Filesystem watcher error", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			ignore := strings.HasSuffix(e.Name, "swp") || strings.HasSuffix(e.Name, "swx")
			if !ignore {
				log.Info("File system event",
					log.String("op", e.Op.String()),
					log.String("name", e.Name))
				if e.Op.Has(fsnotify.Create) {
					// naive way of waiting until the file has finished writing before we get the
					// hash and byte size. If you need to write a file that takes longer than this
					// to finish, you should stop the catalog server. There's no way for us to know
					// if there are no more pending writes
					time.Sleep(1 * time.Second)
					info, err := os.Stat(e.Name)
					if err != nil {
						continue
					}
					_ = srv.library.AddFile(e.Name, info.Size())
				}
				if e.Op.Has(fsnotify.Rename) {
					// this .Name should be the previous name, which is no longer watched
					// the new one should show up in a create (?)
					_ = srv.library.RemoveFile(e.Name)
				}
				if e.Op.Has(fsnotify.Remove) {
					_ = srv.library.RemoveFile(e.Name)
				}
			}
		}
	}
}
