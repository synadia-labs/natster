package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/ConnectEverything/natster/internal/catalogserver"
	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/choria-io/fisk"
)

func ShareCatalog(ctx *fisk.ParseContext) error {
	lib, err := medialibrary.Load(ShareOpts.Name)
	if err != nil {
		return err
	}
	if len(ShareOpts.AccountKey) != 56 ||
		!strings.HasPrefix(ShareOpts.AccountKey, "A") {
		return errors.New("target is not a properly formed account public key")
	}
	lib.AddShare(ShareOpts.AccountKey)
	err = lib.Save()
	if err != nil {
		return err
	}
	fmt.Printf("Shared catalog '%s' with target '%s'. Note: Natster makes no guarantees that the target account exists.\n",
		ShareOpts.Name,
		ShareOpts.AccountKey,
	)

	return nil
}

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
