package catalogserver

import (
	"fmt"
	log "log/slog"
	"slices"
	"strings"

	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/nats.go"
)

// API functions exposed over NATS, topic import/export is used to allow
// the sharing of this API to others

func (srv *CatalogServer) startApiSubscriptions() error {

	_, err := srv.nc.Subscribe(
		fmt.Sprintf("*.natster.catalog.%s.get", srv.library.Name),
		handleCatalogGet(srv))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog get",
			log.String("library", srv.library.Name),
		)
		return err
	}

	return nil
}

func handleCatalogGet(srv *CatalogServer) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		tokens := strings.Split(m.Subject, ".")
		if !slices.Contains(srv.library.Shares, tokens[0]) { // is this account on the sharing list?
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
			Name:    srv.library.Name,
			Entries: convertEntries(catalog),
		}
		catalogRaw := models.NewApiResultPass(catalogSummary)
		if err != nil {
			log.Error(
				"Failed to serialize the catalog",
				log.String("error", err.Error()),
			)
			return
		}

		m.Respond(catalogRaw)
	}
}

func convertEntries(entries []medialibrary.MediaEntry) []models.CatalogEntry {
	out := make([]models.CatalogEntry, len(entries))
	for i, entry := range entries {
		outEntry := models.CatalogEntry{
			Path:        entry.Path,
			Description: entry.Description,
			MimeType:    entry.MimeType,
			ByteSize:    entry.ByteSize,
		}
		out[i] = outEntry
	}
	return out
}
