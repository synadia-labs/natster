package globalservice

import (
	"encoding/json"
	"log/slog"

	"github.com/jellydator/ttlcache/v3"
	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

func handleHeartbeat(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		accountKey := extractAccountKey(m.Subject)
		slog.Info("Receiving heartbeat", slog.String("account", accountKey))
		var hb models.Heartbeat
		err := json.Unmarshal(m.Data, &hb)
		if err != nil {
			slog.Error("Failed to deserialize heartbeat", err)
			return
		}

		srv.hbCache.Set(hb.Catalog, hb, ttlcache.DefaultTTL)
	}
}

func (srv *GlobalService) IsCatalogOnline(catalog string) bool {
	return srv.hbCache.Has(catalog)
}
