package globalservice

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
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

func handleGetContext(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		accountKey := extractAccountKey(m.Subject)
		if accountKey != natsterDotIoAccount {
			slog.Error("Attempt to query context from invalid account",
				slog.String("account", accountKey),
			)
			_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
			return
		}

		nctx, err := srv.GetBoundContextByOAuth(string(m.Data))
		if err != nil {
			slog.Error("Failed to query bound context for OAuth ID", err,
				slog.String("id", string(m.Data)),
			)
		}

		client := syncp.NewAPIClient(syncp.NewConfiguration())
		ctxx := context.WithValue(context.Background(), syncp.ContextServerVariables, map[string]string{
			"baseUrl": "https://cloud.synadia.com",
		})
		ctxx = context.WithValue(ctxx, syncp.ContextAccessToken, nctx.Token)

		creds, _, err := client.NatsUserAPI.DownloadNatsUserCreds(ctxx, nctx.UserID).Execute()
		if err != nil {
			slog.Error("Failed to download NATS user creds from Synadia Cloud", err,
				slog.String("user_id", nctx.UserID),
			)
			return
		}

		resp := models.ContextQueryResponse{
			Context:   *nctx,
			FullCreds: creds,
		}
		_ = m.Respond(models.NewApiResultPass(resp))
	}
}
