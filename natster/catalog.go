package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/choria-io/fisk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
	"github.com/synadia-labs/natster/internal/catalogserver"
	"github.com/synadia-labs/natster/internal/globalservice"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	ngsUrl = "tls://connect.ngs.global"
)

func DescribeCatalogItem(ctx *fisk.ParseContext) error {
	library, err := medialibrary.Load(ShareOpts.Name)
	if err != nil {
		fmt.Printf("Failed to load catalog: %s\n", err)
		return err
	}
	err = library.DescribeItem(DlOpts.Hash, HubOpts.Description)
	if err != nil {
		fmt.Printf("Failed to describe file item in catalog: %s\n", err)
		return err
	}
	err = library.Save()
	if err != nil {
		fmt.Printf("Failed to write updated catalog: %s\n", err)
		return err
	}
	fmt.Printf("New description set for file %s in catalog %s\n", DlOpts.Hash, ShareOpts.Name)

	return nil
}

func ViewCatalogItems(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	client, err := catalogserver.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return err
	}
	items, err := client.GetCatalogItems(ShareOpts.Name)
	if err != nil {
		return err
	}

	t := newTableWriter(fmt.Sprintf("Items in Catalog %s", ShareOpts.Name), "cyan")
	w := t.writer
	w.AppendHeader(table.Row{"Path", "Hash", "Mime Type"})
	for _, item := range items {
		w.AppendRow(table.Row{item.Path, item.Hash, item.MimeType})
	}
	fmt.Println(w.Render())
	return nil
}

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

	catalogs := make(map[string][]models.CatalogShareSummary)
	for _, share := range catshares {
		cat, ok := catalogs[share.Catalog]
		var catshares []models.CatalogShareSummary
		if !ok {
			catshares = make([]models.CatalogShareSummary, 0)
			catshares = append(catshares, share)
		} else {
			catshares = cat
			catshares = append(catshares, share)
		}
		catalogs[share.Catalog] = catshares
	}

	t := newTableWriter("Shared Catalogs", "cyan")
	w := t.writer
	w.AppendHeader(table.Row{"", "Catalog", "From", "To"})

	for catalog, shares := range catalogs {
		online := " "
		if shares[0].CatalogOnline {
			online = "🟢"
		}
		for i, share := range shares {
			if share.FromAccount == nctx.AccountPublicKey {
				share.FromAccount = "me"
			}
			if share.ToAccount == nctx.AccountPublicKey {
				share.ToAccount = "me"
			}
			if i > 0 {
				w.AppendRow(table.Row{"", "", share.FromAccount, share.ToAccount})
			} else {
				w.AppendRow(table.Row{online, catalog, share.FromAccount, share.ToAccount})
			}
		}

	}
	fmt.Println(w.Render())

	return nil
}

type inboxResponse struct {
	UnimportedShares []models.CatalogShareSummary `json:"unimported_shares"`
}

func RenderInbox(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	ctxOpts := []natscontext.Option{
		natscontext.WithServerURL(ngsUrl),
		natscontext.WithCreds(nctx.CredsPath),
	}

	natsContext, err := natscontext.New("natster_temp", false, ctxOpts...)
	if err != nil {
		return err
	}

	conn, err := natsContext.Connect(nats.Name("natster_client"))
	if err != nil {
		return err
	}

	resp, err := conn.Request("natster.local.inbox", []byte{}, 1*time.Second)
	if err != nil {
		return err
	}
	var inbox models.TypedApiResult[inboxResponse]
	err = json.Unmarshal(resp.Data, &inbox)
	if err != nil {
		return err
	}
	if inbox.Code != 200 {
		fmt.Printf("Error retrieving inbox: %s", *inbox.Error)
		return nil
	}
	if len(inbox.Data.UnimportedShares) == 0 {
		fmt.Println("Inbox zero! No unimported catalogs")
		return nil
	}

	catalogNames := make([]string, len(inbox.Data.UnimportedShares))
	for i, cat := range inbox.Data.UnimportedShares {
		catalogNames[i] = cat.Catalog
	}

	userPrompt := survey.Select{
		Message: "Select a catalog to import",
		Options: catalogNames,
	}
	var selectedCatalog int
	err = survey.AskOne(&userPrompt, &selectedCatalog, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	ShareOpts.Name = catalogNames[selectedCatalog]
	ShareOpts.AccountKey = inbox.Data.UnimportedShares[selectedCatalog].FromAccount

	return ImportCatalog(ctx)
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
	mediaImportName := fmt.Sprintf("natster_%s_media", ShareOpts.Name)

	imports, _, err := client.AccountAPI.ListSubjectImports(ctxx, nctx.AccountID).Execute()
	if err != nil {
		fmt.Printf("Failed to query subject imports for this account: %s\n", err)
		return err
	}
	catFound := false
	mediaFound := false

	for _, imp := range imports.Items {
		if *imp.JwtSettings.Name == importName {
			fmt.Printf("✅ Catalog '%s' is imported from account '%s'\n",
				ShareOpts.Name, ShareOpts.AccountKey,
			)
			catFound = true
		}
		if *imp.JwtSettings.Name == mediaImportName {
			fmt.Printf("✅ Catalog '%s' media stream is imported from account '%s'\n",
				ShareOpts.Name, ShareOpts.AccountKey,
			)
			mediaFound = true
		}
	}
	if !catFound {
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
		fmt.Printf("✅ Catalog '%s' is imported from account '%s'\n",
			ShareOpts.Name, ShareOpts.AccountKey,
		)
	}
	if !mediaFound {
		importReq := syncp.SubjectImportCreateRequest{
			JwtSettings: syncp.Import{
				Account:      syncp.Ptr(ShareOpts.AccountKey),
				Subject:      syncp.Ptr(fmt.Sprintf("%s.natster.media.%s.*", nctx.AccountPublicKey, ShareOpts.Name)),
				LocalSubject: syncp.Ptr(fmt.Sprintf("natster.media.%s.*", ShareOpts.Name)),
				Name:         syncp.Ptr(mediaImportName),
				Type:         syncp.Ptr(syncp.EXPORTTYPE_STREAM),
			},
		}
		_, _, err = client.AccountAPI.CreateSubjectImport(ctxx, nctx.AccountID).SubjectImportCreateRequest(importReq).Execute()
		if err != nil {
			return err
		}
		fmt.Printf("✅ Catalog '%s' media stream is imported from account '%s'\n",
			ShareOpts.Name, ShareOpts.AccountKey,
		)
	}
	return nil
}

func UnshareCatalog(ctx *fisk.ParseContext) error {
	if len(ShareOpts.AccountKey) != 56 ||
		!strings.HasPrefix(ShareOpts.AccountKey, "A") {
		return errors.New("target is not a properly formed account public key")
	}

	ShareOpts.Name = strings.ToLower(ShareOpts.Name)
	err := publishCatalogUnshared()
	if err != nil {
		return err
	}

	fmt.Printf("Stopped sharing catalog '%s' with target '%s'. Note: Natster makes no guarantees that the target account exists.\n",
		ShareOpts.Name,
		ShareOpts.AccountKey,
	)
	fmt.Println("If this catalog wasn't shared with the target or has already been unshared, this operation will have no effect.")

	return nil
}

func ShareCatalog(ctx *fisk.ParseContext) error {

	if len(ShareOpts.AccountKey) != 56 ||
		!strings.HasPrefix(ShareOpts.AccountKey, "A") {
		return errors.New("target is not a properly formed account public key")
	}

	nctx, err := loadContext()
	if err != nil {
		return err
	}
	if nctx.AccountPublicKey == ShareOpts.AccountKey {
		fmt.Println("You cannot share catalogs with yourself.")
		return nil
	}

	client, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		slog.Error(
			"Failed to connect to NGS",
			slog.String("error", err.Error()),
		)
		return err
	}
	validationResult, err := client.ValidateCatalogName(ShareOpts.Name)
	if err != nil {
		return err
	}
	if !validationResult.Valid {
		fmt.Printf("The catalog named '%s' cannot be shared: %s\n", ShareOpts.Name, validationResult.Message)
		return nil
	}

	ShareOpts.Name = strings.ToLower(ShareOpts.Name)
	err = publishCatalogShared()
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

	nctx, _ := loadContext()
	client, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	// skip this check if we can't talk to NGS
	if err == nil {
		validationResult, err := client.ValidateCatalogName(HubOpts.Name)
		if err == nil {
			if !validationResult.Valid {
				fmt.Printf("The new catalog '%s' failed validation, and will not be shareable outside your account: %s\n. Catalog not created.\n",
					HubOpts.Name,
					validationResult.Message,
				)
				return nil
			}
		}
	}

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

func publishCatalogUnshared() error {
	ctx, _ := loadContext()
	client, err := globalservice.NewClientWithCredsPath(ctx.CredsPath)
	if err != nil {
		slog.Error(
			"Failed to connect to NGS",
			slog.String("error", err.Error()),
		)
		return err
	}
	err = client.PublishEvent(models.CatalogUnsharedEventType, ShareOpts.Name, ShareOpts.AccountKey, nil)
	if err != nil {
		slog.Error(
			"Failed to publish catalog unshared event",
			slog.String("error", err.Error()),
		)
		return err
	}
	client.Drain()

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
	client.Drain()

	return nil
}
