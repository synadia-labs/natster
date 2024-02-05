package globalservice

import (
	"encoding/json"
	"log/slog"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/jellydator/ttlcache/v3"
	"github.com/nats-io/nats.go"
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
		srv.hbCache.Set(accountKey, hb, ttlcache.DefaultTTL)
	}
}
