package hubserver

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

func (hub *HubServer) startApiSubscriptions() error {

	_, err := hub.nc.Subscribe(
		fmt.Sprintf("%s.catalog.%s.get", APIPrefix, hub.library.Name),
		handleCatalogGet(hub))
	if err != nil {
		log.Error(
			"Failed to subscribe to catalog get",
			"library", hub.library.Name,
		)
		return err
	}

	return nil
}

func handleCatalogGet(hub *HubServer) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		catalog, err := hub.library.GetCatalog()
		if err != nil {
			log.Error(
				"Failed to query the library catalog",
				"error", err,
			)
			return
		}
		fmt.Printf("%+v\n", catalog)
		catalogSummary := models.CatalogSummary{
			Name:    hub.library.Name,
			Entries: convertEntries(catalog),
		}
		catalogRaw, err := json.Marshal(catalogSummary)
		if err != nil {
			log.Error(
				"Failed to serialize the catalog",
				"error", err,
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
