package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/ConnectEverything/natster/internal/catalogserver"
	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/choria-io/fisk"
)

func StartCatalogServer(ctx *fisk.ParseContext) error {
	ctxx := context.Background()

	nctx, err := loadContext()
	if err != nil {
		return err
	}
	// TODO: clean this up, scar tissue
	Opts.Creds = nctx.CredsPath
	catalogserver.CliStart(Opts, HubOpts)
	<-ctxx.Done()

	return nil
}

func NewCatalog(ctx *fisk.ParseContext) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, nil)))

	lib, err := medialibrary.New(HubOpts.RootPath, HubOpts.Name, HubOpts.Description)
	if err != nil {
		fmt.Printf("Failed to initialize media catalog: %s\n", err)
		return err
	}
	err = lib.Ingest()
	if err != nil {
		fmt.Printf("Failed to read in media files: %s\n", err)
		return err
	}
	err = lib.Save()
	if err != nil {
		fmt.Printf("Failed to store new catalog: %s\n", err)
		return err
	}
	fmt.Printf("New catalog created: %s\n", HubOpts.Name)
	return nil
}
