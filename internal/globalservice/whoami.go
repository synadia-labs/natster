package globalservice

import (
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

func handleWhoAmi(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		accountKey := extractAccountKey(m.Subject)
		oauth, err := srv.GetOAuthIdForAccount(accountKey)
		if err != nil {
			slog.Error("Failed to query OAuth ID for account", err)
			_ = m.Respond(models.NewApiResultFail("Not Found", 404))
			return
		}
		// Note: a non-error but nil oauth is valid - just means it hasn't been context
		// bound yet

		resp := models.WhoamiResponse{
			AccountKey:    accountKey,
			OAuthIdentity: oauth,
		}
		_ = m.Respond(models.NewApiResultPass(resp))
	}
}
