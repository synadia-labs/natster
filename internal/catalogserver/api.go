package catalogserver

import (
	"encoding/json"
	"fmt"
	log "log/slog"

	"github.com/ConnectEverything/natster/internal/medialibrary"
	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/nats.go"
)

// API functions exposed over NATS, topic import/export is used to allow
// the sharing of this API to others

const (
	APIPrefix = "natster"
)

func (srv *CatalogServer) startApiSubscriptions() error {

	_, err := srv.nc.Subscribe(
		fmt.Sprintf("%s.catalog.%s.get", APIPrefix, srv.library.Name),
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
		catalogRaw, err := json.Marshal(catalogSummary)
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
