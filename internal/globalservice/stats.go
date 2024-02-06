package globalservice

import (
	"encoding/json"
	"log/slog"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/nats.go"
)

func handleStats(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		initialized, err := srv.GetTotalInitializedAccounts()
		if err != nil {
			return
		}
		stats := models.CommunityStats{
			TotalInitialized: initialized,
			RunningCatalogs:  uint64(srv.hbCache.Len()),
		}
		bytes, err := json.Marshal(&stats)
		if err != nil {
			slog.Error("Failed to serialize community stats", err)
			return
		}
		_ = m.Respond(bytes)
	}
}
