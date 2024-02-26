package catalogserver

import (
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

func CliStart(ctx *models.NatsterContext, opts *models.Options, hubopts *models.HubOptions) {
	nc, err := generateConnectionFromOpts(opts)
	if err != nil {
		log.Error(
			"Failed to connect to NATS",
			log.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info(
		"Established Natster NATS connection",
		log.String("servers", synadiaCloudServer),
	)

	library, err := medialibrary.Load(hubopts.Name)
	if err != nil {
		log.Error(
			"Failed to open media catalog",
			log.String("name", hubopts.Name),
		)
		os.Exit(1)
	}
	log.Info(
		"Opened Media Catalog",
		log.String("name", hubopts.Name),
		log.String("rootpath", library.RootDir),
	)

	server := New(ctx, nc, library, hubopts.AllowAll)
	err = server.Start()
	if err != nil {
		log.Error(
			"Failed to start Natster Hub",
			log.String("error", err.Error()),
		)
		os.Exit(1)
	}

	setupSignalHandlers(server)
}

func setupSignalHandlers(hub *CatalogServer) {
	go func() {
		signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				log.Info(
					"Caught signal, requesting clean shutdown",
					log.String("signal", s.String()),
				)
				err := hub.Stop()
				if err != nil {
					log.Error(
						"Hub server failed to stop",
						log.String("error", err.Error()),
					)
				}
				os.Exit(0)
			case s == syscall.SIGQUIT:
				log.Info(
					"Caught quit signal, still trying graceful shutdown",
					log.String("signal", s.String()),
				)
				err := hub.Stop()
				if err != nil {
					log.Error(
						"Hub server failed to stop",
						log.String("error", err.Error()),
					)
				}
				os.Exit(0)
			}
		}
	}()
}
