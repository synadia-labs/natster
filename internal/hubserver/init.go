package hubserver

import (
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/ConnectEverything/natster/internal/models"
)

func CliStart(opts *models.Options, hubopts *models.HubOptions) {
	nc, err := generateConnectionFromOpts(opts)
	if err != nil {
		log.Error(
			"Failed to connect to NATS",
			log.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info(
		"Established Natster Hub NATS connection",
		log.String("servers", opts.Servers),
	)

	library, err := medialibrary.New(nc, hubopts.RootPath, hubopts.Name, hubopts.Description)
	if err != nil {
		log.Error(
			"Failed to create media library",
			log.String("path", hubopts.RootPath),
			log.String("name", hubopts.Name),
		)
		os.Exit(1)
	}
	log.Info(
		"Opened Media Library",
		log.String("path", hubopts.RootPath),
		log.String("name", hubopts.Name),
	)

	server := New(nc, library)
	err = server.Start(hubopts.Port)
	if err != nil {
		log.Error(
			"Failed to start Natster Hub",
			log.String("error", err.Error()),
		)
		os.Exit(1)
	}

	setupSignalHandlers(server)
}

func setupSignalHandlers(hub *HubServer) {
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
