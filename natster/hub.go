package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ConnectEverything/natster/internal/hubserver"
	"github.com/choria-io/fisk"
)

func HubUp(ctx *fisk.ParseContext) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	ctxx := context.Background()
	hubserver.CliStart(Opts, HubOpts)
	<-ctxx.Done()

	return nil
}
