package globalservice

import (
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

func handleMyShares(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		key := extractAccountKey(m.Subject)
		shares, err := srv.GetMyCatalogs(key)
		if err != nil {
			slog.Error("Failed to handle request for my catalog shares", slog.Any("error", err))
			return
		}
		res := models.NewApiResultPass(shares)
		_ = m.Respond(res)
	}
}

func handleStats(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		initialized, err := srv.GetTotalInitializedAccounts()
		if err != nil {
			return
		}

		// TODO: replace this calculation with a global stats projection so we
		// can just query it
		shareCount, err := srv.GetTotalSharedCatalogs()
		if err != nil {
			return
		}
		stats := models.CommunityStats{
			TotalInitialized: initialized,
			RunningCatalogs:  uint64(srv.hbCache.Len()),
			SharedCatalogs:   shareCount,
		}
		res := models.NewApiResultPass(stats)
		_ = m.Respond(res)
	}
}
