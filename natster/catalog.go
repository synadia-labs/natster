package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
	"github.com/synadia-labs/natster/internal/catalogserver"
	"github.com/synadia-labs/natster/internal/globalservice"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

func ListCatalogs(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	client, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return err
	}
	catshares, err := client.GetMyCatalogs()
	if err != nil {
		return err
	}
	for _, share := range catshares {
		if share.FromAccount == nctx.AccountPublicKey {
			fmt.Printf("* ")
		}
		fmt.Println(share.Catalog)
	}

	return nil
}

func ImportCatalog(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	client := syncp.NewAPIClient(syncp.NewConfiguration())
	ctxx := context.WithValue(context.Background(), syncp.ContextServerVariables, map[string]string{
		"baseUrl": baseUrl,
	})
	ctxx = context.WithValue(ctxx, syncp.ContextAccessToken, nctx.Token)

	importName := fmt.Sprintf("natster_%s", ShareOpts.Name)

	imports, _, err := client.AccountAPI.ListSubjectImports(ctxx, nctx.AccountID).Execute()
	if err != nil {
		fmt.Printf("Failed to query subject imports for this account: %s\n", err)
		return err
	}
	for _, imp := range imports.Items {
		if *imp.JwtSettings.Name == importName {
			fmt.Printf("✅ Catalog '%s' is already imported from account '%s'\n",
				ShareOpts.Name, ShareOpts.AccountKey,
			)
			return nil
		}
	}

	importReq := syncp.SubjectImportCreateRequest{
		JwtSettings: syncp.Import{
			Account:      syncp.Ptr(ShareOpts.AccountKey),
			Subject:      syncp.Ptr(fmt.Sprintf("%s.natster.catalog.%s.>", nctx.AccountPublicKey, ShareOpts.Name)),
			LocalSubject: syncp.Ptr(fmt.Sprintf("natster.catalog.%s.>", ShareOpts.Name)),
			Name:         syncp.Ptr(importName),
			Type:         syncp.Ptr(syncp.EXPORTTYPE_SERVICE),
		},
	}
	_, _, err = client.AccountAPI.CreateSubjectImport(ctxx, nctx.AccountID).SubjectImportCreateRequest(importReq).Execute()
	if err != nil {
		return err
	}

	fmt.Printf("Catalog '%s' imported from account '%s'. You can now query this catalog.\n", ShareOpts.Name, ShareOpts.AccountKey)
	return nil
}

func ShareCatalog(ctx *fisk.ParseContext) error {

	if len(ShareOpts.AccountKey) != 56 ||
		!strings.HasPrefix(ShareOpts.AccountKey, "A") {
		return errors.New("target is not a properly formed account public key")
	}

	ShareOpts.Name = strings.ToLower(ShareOpts.Name)
	err := publishCatalogShared()
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
	catalogserver.CliStart(nctx, Opts, HubOpts)
	<-ctxx.Done()

	return nil
}

func NewCatalog(ctx *fisk.ParseContext) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, nil)))

	HubOpts.Name = strings.ToLower(HubOpts.Name)
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

func publishCatalogShared() error {
	ctx, _ := loadContext()
	client, err := globalservice.NewClientWithCredsPath(ctx.CredsPath)
	if err != nil {
		slog.Error(
			"Failed to connect to NGS",
			slog.String("error", err.Error()),
		)
		return err
	}
	err = client.PublishEvent(models.CatalogSharedEventType, ShareOpts.Name, ShareOpts.AccountKey, nil)
	if err != nil {
		slog.Error(
			"Failed to publish catalog shared event",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
