package globalservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	synadiaHubAccount = "AC5V4OC2POUAX4W4H7CKN5TQ5AKVJJ4AJ7XZKNER6P6DHKBYGVGJHSNC"
)

// Returns the total number of accounts in which developers have run `natster init`
func (srv *GlobalService) GetTotalInitializedAccounts() (uint64, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate key value bucket", slog.Any("error", err), slog.String("bucket", accountProjectionBucketName))
		return 0, err
	}
	status, err := kv.Status(ctx)
	if err != nil {
		slog.Error("Couldn't obtain status of key value bucket", slog.Any("error", err), slog.String("bucket", accountProjectionBucketName))
		return 0, err
	}

	// NOTE: this number will only be accurate if we're only keeping the latest version (e.g. no history)
	// in the kv bucket
	return status.Values(), nil
}

// This number is currently inaccurate due to the "unshare" feature. TODO: we can fix this when we add a global
// stats projection
func (srv *GlobalService) GetTotalSharedCatalogs() (uint64, error) {
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.CatalogSharedEventType)
	return srv.countFilteredEvents(subject)
}

// Retrieves the most recent bound context by loading the account projection
func (srv *GlobalService) GetBoundContext(myAccountKey string) (*models.ContextBoundEvent, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate key value bucket", slog.Any("error", err), slog.String("bucket", accountProjectionBucketName))
		return nil, err
	}
	myAccount, err := loadAccount(kv, myAccountKey)
	if err != nil {
		slog.Error("Failed to load source account for catalog query", slog.Any("error", err))
		return nil, err
	}

	return myAccount.BoundContext, nil
}

func (srv *GlobalService) GetOAuthIdForAccount(accountKey string) (*string, error) {
	discoveredContext, err := srv.GetBoundContext(accountKey)
	if err != nil {
		return nil, err
	}

	return &discoveredContext.OAuthIdentity, nil
}

// TODO: this needs to be converted to use a projection rather than running through the
// event stream every time
func (srv *GlobalService) GetBoundContextByOAuth(oauthId string) (*models.NatsterContext, error) {
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.ContextBoundEventType)
	js, err := jetstream.New(srv.nc)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return nil, err
	}
	consumer, err := s.CreateConsumer(context.Background(), jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckNonePolicy,
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan *models.NatsterContext)

	cc, _ := consumer.Consume(func(msg jetstream.Msg) {
		var discoveredContext models.ContextBoundEvent
		err := json.Unmarshal(msg.Data(), &discoveredContext)
		if err != nil {
			slog.Error("Deserialization failure of context bound event", err)
			ch <- nil
		}
		if discoveredContext.OAuthIdentity == oauthId {
			ch <- &discoveredContext.BoundContext
		}
	})
	select {
	case discoveredContext := <-ch:
		cc.Stop()
		return discoveredContext, nil
	case <-time.After(1300 * time.Millisecond):
		cc.Stop()
		return nil, nil
	}
}

// Reads the account projection for the given querying account and returns a flattened list
// of directional shares to and from this account
func (srv *GlobalService) GetMyCatalogs(myAccountKey string) ([]models.CatalogShareSummary, error) {

	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate key value bucket", slog.Any("error", err), slog.String("bucket", accountProjectionBucketName))
		return nil, err
	}
	myAccount, err := loadAccount(kv, myAccountKey)
	if err != nil {
		slog.Error("Failed to load source account for catalog query", slog.Any("error", err))
		return nil, err
	}
	summaries := make([]models.CatalogShareSummary, 0)
	for _, outShare := range myAccount.OutShares {
		summaries = append(summaries, models.CatalogShareSummary{
			FromAccount:   myAccountKey,
			ToAccount:     outShare.Account,
			Catalog:       outShare.Catalog,
			CatalogOnline: srv.IsCatalogOnline(outShare.Catalog),
			Revision:      srv.CatalogRevision(outShare.Catalog),
		})
	}
	for _, inShare := range myAccount.InShares {
		summaries = append(summaries, models.CatalogShareSummary{
			FromAccount:   inShare.Account,
			ToAccount:     myAccountKey,
			Catalog:       inShare.Catalog,
			CatalogOnline: srv.IsCatalogOnline(inShare.Catalog),
			Revision:      srv.CatalogRevision(inShare.Catalog),
		})
	}

	return summaries, nil
}

func (srv *GlobalService) countFilteredEvents(subject string) (uint64, error) {
	js, err := jetstream.New(srv.nc)
	if err != nil {
		return 0, err
	}
	ctx := context.Background()
	opts := make([]jetstream.StreamInfoOpt, 0)
	opts = append(opts, jetstream.WithSubjectFilter(subject))
	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return 0, err
	}
	streamInfo, err := s.Info(ctx, opts...)
	if err != nil {
		return 0, err
	}
	count := uint64(0)
	for _, v := range streamInfo.State.Subjects {
		count += v
	}
	return count, nil
}

func handleEventPut(srv *GlobalService) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		key := extractAccountKey(m.Subject)
		var evt models.NatsterEvent
		err := json.Unmarshal(m.Data, &evt)
		if err != nil {
			slog.Error("Failed to deserialize Natster event", err)
			return
		}

		// write the event to the stream
		subject := fmt.Sprintf("natster.events.%s.%s.%s.%s", key, evt.Target, evt.Catalog, evt.EventType)
		raw, err := json.Marshal(evt.Data)
		if err != nil {
			slog.Error("Failed to serialize Natster event", err)
			return
		}
		slog.Info("Writing Natster global event",
			slog.Int("bytes", len(raw)),
			slog.String("target", evt.Target),
			slog.String("catalog", evt.Catalog),
			slog.String("event_type", evt.EventType),
		)
		err = srv.nc.Publish(subject, raw)
		if err != nil {
			slog.Error("Failed to publish Natster event", err)
			return
		}

		// NOTE: in a real event sourced system, we'd have a consumer listening to this event
		// type and when we receive it, publish the autoshare and then ack. But that's an exercise
		// we can leave for when Natster has a hojillion users
		if evt.EventType == models.NatsterInitializedEventType {
			err = srv.publishSynadiaHubAutoShare(key)
			if err != nil {
				slog.Error("Failed to publish autoshare event for synadia hub", slog.Any("error", err))
			}
		}
	}
}

// NOTE: this is safe to publish now because we're no longer writing multiple initialized events
func (srv *GlobalService) publishSynadiaHubAutoShare(targetKey string) error {
	slog.Info("Detected Natster initialized event, auto-sharing synadia hub.")
	subject := fmt.Sprintf("natster.events.%s.%s.synadiahub.%s",
		synadiaHubAccount,
		targetKey,
		models.CatalogSharedEventType)

	// This is intentional - the catalog shared event is empty (for now) as all the relevant
	// data is in the subject tokens
	err := srv.nc.Publish(subject, []byte{})
	if err != nil {
		return err
	}

	return nil
}
