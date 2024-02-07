package main

import (
	"context"
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/globalservice"
)

func main() {
	jwt := os.Getenv("USER_JWT")
	seed := os.Getenv("USER_SEED")
	server := "tls://connect.ngs.global"

	nc, err := nats.Connect(server,
		nats.UserJWTAndSeed(jwt, seed),
		nats.Name("natster_global"),
	)
	if err != nil {
		panic(err)
	}

	ctxx := context.Background()
	srv := globalservice.New(nc)
	err = srv.Start()
	if err != nil {
		panic(err)
	}
	setupSignalHandlers(srv)

	<-ctxx.Done()

}

func setupSignalHandlers(hub *globalservice.GlobalService) {
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
						"Global service server failed to stop",
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
						"Global service server failed to stop",
						log.String("error", err.Error()),
					)
				}
				os.Exit(0)
			}
		}
	}()
}
