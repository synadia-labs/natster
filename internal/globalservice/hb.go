package globalservice

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jellydator/ttlcache/v3"
	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

func handleHeartbeat(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		accountKey := extractAccountKey(m.Subject)
		slog.Debug("Receiving heartbeat", slog.String("account", accountKey))
		var hb models.Heartbeat
		err := json.Unmarshal(m.Data, &hb)
		if err != nil {
			slog.Error("Failed to deserialize heartbeat", err)
			return
		}

		srv.hbCache.Set(hb.Catalog, hb, ttlcache.DefaultTTL)

		go srv.rebroadcastHeartbeat(accountKey, m.Data)
	}
}

func (srv *GlobalService) IsCatalogOnline(catalog string) bool {
	return srv.hbCache.Has(catalog)
}

func (srv *GlobalService) CatalogRevision(catalog string) int64 {
	c := srv.hbCache.Get(catalog)
	if c == nil {
		return 0
	}
	return c.Value().Revision
}

func (srv *GlobalService) rebroadcastHeartbeat(accountKey string, data []byte) {
	shares, err := srv.GetMyCatalogs(accountKey)
	if err != nil {
		slog.Error("Failed to get catalog list for source account",
			slog.String("account", accountKey))
	}

	// Send a heartbeat to each account _to which_ the catalog has been shared from
	// the source account (accountKey)
	for _, share := range shares {
		if share.FromAccount == accountKey {
			subject := fmt.Sprintf("%s.natster.global-events.heartbeat", share.ToAccount)
			_ = srv.nc.Publish(subject, data)
		}
	}
}
