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
			return
		}
		res := models.NewApiResultPass(shares)

		if err != nil {
			slog.Error("Failed to serialize share summaries")
			return
		}
		_ = m.Respond(res)
	}
}

func handleStats(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		initialized, err := srv.GetTotalInitializedAccounts()
		if err != nil {
			return
		}
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
		if err != nil {
			slog.Error("Failed to serialize community stats", err)
			return
		}
		_ = m.Respond(res)
	}
}
