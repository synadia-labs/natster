package catalogserver

import (
	"fmt"
	log "log/slog"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

// API functions exposed over NATS, topic import/export is used to allow
// the sharing of this API to others

func (srv *CatalogServer) startApiSubscriptions() error {

	_, err := srv.nc.QueueSubscribe(
		fmt.Sprintf("*.natster.catalog.%s.get", srv.library.Name), "natstercatalog",
		handleCatalogGet(srv, false))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog get",
			log.String("library", srv.library.Name),
		)
		return err
	}

	_, err = srv.nc.QueueSubscribe(
		fmt.Sprintf("natster.catalog.%s.get", srv.library.Name), "natstercatalog",
		handleCatalogGet(srv, true))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog get",
			log.String("library", srv.library.Name),
		)
		return err
	}

	_, err = srv.nc.QueueSubscribe(
		fmt.Sprintf("*.natster.catalog.%s.download", srv.library.Name), "natstercatalog",
		handleDownloadRequest(srv, false))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog item download",
			log.String("library", srv.library.Name),
		)
		return err
	}

	_, err = srv.nc.QueueSubscribe(
		fmt.Sprintf("natster.catalog.%s.download", srv.library.Name), "natstercatalog",
		handleDownloadRequest(srv, true))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog item download",
			log.String("library", srv.library.Name),
		)
		return err
	}

	_, err = srv.nc.QueueSubscribe("natster.local.>", "natsterlocalservices",
		handleLocalServiceRequest(srv))
	if err != nil {
		log.Error(
			"Failed to subscribe to natster local services subject")
		return err
	}

	return nil
}

func (srv *CatalogServer) isClientAllowed(accountKey string) bool {
	if srv.allowAll {
		return true
	}
	cats, err := srv.globalServiceClient.GetMyCatalogs()
	if err != nil {
		return false
	}
	if cats == nil {
		return false
	}
	// Is there a sharing record from me to the calling account
	// for the catalog in question?
	allowed := false
	for _, cat := range *cats {
		if cat.Catalog == srv.library.Name &&
			cat.ToAccount == accountKey {
			allowed = true
		}
	}
	return allowed
}

func handleCatalogGet(srv *CatalogServer, local bool) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		tokens := strings.Split(m.Subject, ".")

		allowed := true
		if !local {
			allowed = srv.isClientAllowed(tokens[0])
		}

		if !allowed { // is this account on the sharing list?
			_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
			return
		}

		catalog, err := srv.library.GetCatalog()
		if err != nil {
			log.Error(
				"Failed to query the library catalog",
				log.String("error", err.Error()),
			)
			return
		}
		catalogSummary := models.CatalogSummary{
			Name:        srv.library.Name,
			Description: srv.library.Description,
			Image:       srv.library.ImageLink,
			Entries:     convertEntries(catalog),
		}
		catalogRaw := models.NewApiResultPass(catalogSummary)
		m.Respond(catalogRaw)
	}
}

func convertEntries(entries []*medialibrary.MediaEntry) []models.CatalogEntry {
	out := make([]models.CatalogEntry, len(entries))
	for i, entry := range entries {
		outEntry := models.CatalogEntry{
			Path:        entry.Path,
			Hash:        entry.Hash,
			Description: entry.Description,
			MimeType:    entry.MimeType,
			ByteSize:    entry.ByteSize,
		}
		out[i] = outEntry
	}
	return out
}
