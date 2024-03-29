package main

import (
	"fmt"

	"github.com/choria-io/fisk"
	"github.com/synadia-labs/natster/internal/catalogserver"
)

func DownloadFile(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	catClient, err := catalogserver.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return err
	}

	err = catClient.DownloadFile(ShareOpts.Name, DlOpts.Hash, DlOpts.OutputPath)
	if err != nil {
		fmt.Printf("Failed to download file: %s\n", err)
		return err
	}

	return nil

}
