package globalservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	shorturl "github.com/aviddiviner/shortcode-go"
	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	tokenValidTimeMinutes = 5
	natsterDotIoAccount   = "AA2JVG74M2LCCNYYFMBAANHRNFTUAVWZJGTDAREKDBE23DRLBRWYQNLD"
)

func handleOtcClaim(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		var request models.OtcClaimRequest
		accountKey := extractAccountKey(m.Subject)
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			_ = m.Respond(models.NewApiResultFail("Bad Request", 400))
			return
		}
		if accountKey != natsterDotIoAccount {
			slog.Error("Attempt to claim OTC from non-website account")
			_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
			return
		}

		contextRequest, err := srv.readOneTimeCode(request.Code)
		if err != nil {
			slog.Error("Failed to read code from bucket. Possibly expired?", err)
			_ = m.Respond(models.NewApiResultFail("Not Found", 404))
			return
		}

		evtData := models.ContextBoundEvent{
			OAuthIdentity: request.OAuthIdentity,
			BoundContext:  *contextRequest,
		}

		// write the event to the stream
		// natster.events.{key}.{target}.{catalog}.{event_type}
		subject := fmt.Sprintf("natster.events.%s.none.none.%s", contextRequest.AccountPublicKey, models.ContextBoundEventType)
		raw, err := json.Marshal(evtData)
		if err != nil {
			slog.Error("Failed to serialize Natster event", err)
			return
		}
		slog.Info("Writing Natster global event",
			slog.Int("bytes", len(raw)),
			slog.String("target", "none"),
			slog.String("catalog", "none"),
			slog.String("event_type", models.ContextBoundEventType),
		)
		err = srv.nc.Publish(subject, raw)
		if err != nil {
			slog.Error("Failed to publish Natster event", err)
			return
		}

		_ = m.Respond(models.NewApiResultPass(evtData.BoundContext))
	}
}

func handleOtcGenerate(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		var request models.NatsterContext
		err := json.Unmarshal(m.Data, &request)
		if err != nil {
			slog.Error("Failed to deserialize OTC request", err)
			_ = m.Respond(models.NewApiResultFail("Bad Request", 400))
			return
		}

		r := rand.New(rand.NewSource(time.Now().Unix()))
		otc := shorturl.EncodeID(r.Int())
		accountKey := extractAccountKey(m.Subject)
		if accountKey != request.AccountPublicKey {
			slog.Error("Account trying to generate code for a different context", errors.New("sneaky"),
				slog.String("valid_key", accountKey),
				slog.String("attempted_key", request.AccountPublicKey),
			)
			_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
			return
		}

		resp := models.OtcGenerateResponse{
			Code:         otc,
			ClaimUrl:     fmt.Sprintf("https://natster.io/login/%s", otc),
			ValidMinutes: tokenValidTimeMinutes,
		}

		err = srv.writeOneTimeCode(otc, request)
		if err != nil {
			slog.Error("Failed to write one-time code", err)
			_ = m.Respond(models.NewApiResultFail("Internal Server Error", 500))
			return
		}

		_ = m.Respond(models.NewApiResultPass(resp))
	}
}

func (srv *GlobalService) writeOneTimeCode(code string, request models.NatsterContext) error {
	kv, err := srv.createOrReuseOtcBucket()
	if err != nil {
		return err
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = kv.Put(context.Background(), code, requestBytes)
	if err != nil {
		return err
	}

	return nil
}

func (srv *GlobalService) readOneTimeCode(code string) (*models.NatsterContext, error) {
	kv, err := srv.createOrReuseOtcBucket()
	if err != nil {
		return nil, err
	}
	entry, err := kv.Get(context.Background(), code)
	if err != nil {
		return nil, err
	}
	bytes := entry.Value()
	var originalContext models.NatsterContext
	err = json.Unmarshal(bytes, &originalContext)
	if err != nil {
		return nil, err
	}
	return &originalContext, nil
}
