package globalservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
	"github.com/synadia-labs/natster/internal/models"
)

func handleWhoAmi(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		accountKey := extractAccountKey(m.Subject)
		js, _ := jetstream.New(srv.nc)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		kv, err := js.KeyValue(ctx, accountProjectionBucketName)
		if err != nil {
			slog.Error("Failed to get key value store", slog.Any("error", err))
			_ = m.Respond(models.NewApiResultFail("Internal server error", 500))
			return
		}
		act, _ := loadAccount(kv, accountKey)
		if act == nil {
			_ = m.Respond(models.NewApiResultFail("Not found", 404))
			return
		}

		// Note: a non-error but nil oauth is valid - just means it hasn't been context
		// bound yet
		oauth, err := srv.GetOAuthIdForAccount(accountKey)
		if err != nil {
			slog.Error("Failed to query OAuth ID for account", err)
			_ = m.Respond(models.NewApiResultFail("Internal server error", 500))
			return
		}

		resp := models.WhoamiResponse{
			AccountKey:    accountKey,
			OAuthIdentity: oauth,
			Initialized:   act.InitializedAt,
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
			_ = m.Respond(models.NewApiResultFail("Not Found", 404))
			return
		}
		if nctx == nil {
			// no need to log this, as it'll ultimately be the most common case
			_ = m.Respond(models.NewApiResultFail("Not Found", 404))
			return
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
			_ = m.Respond(models.NewApiResultFail("Not Found", 404))
			return
		}

		jwt, err := nkeys.ParseDecoratedJWT([]byte(creds))
		if err != nil {
			slog.Error("Corrupt credentials data", err)
			_ = m.Respond(models.NewApiResultFail("Internal Server Error", 500))
			return
		}
		seed, err := nkeys.ParseDecoratedUserNKey([]byte(creds))
		if err != nil {
			slog.Error("Corrupt credentials data", err)
			_ = m.Respond(models.NewApiResultFail("Internal Server Error", 500))
		}
		seedKey, _ := seed.Seed()

		resp := models.ContextQueryResponse{
			Context:  *nctx,
			UserJwt:  jwt,
			UserSeed: string(seedKey),
		}
		_ = m.Respond(models.NewApiResultPass(resp))
	}
}
