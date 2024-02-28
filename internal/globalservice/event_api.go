package globalservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.NatsterInitializedEventType)
	return srv.countFilteredEvents(subject)
}

func (srv *GlobalService) GetTotalSharedCatalogs() (uint64, error) {
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.CatalogSharedEventType)
	return srv.countFilteredEvents(subject)
}

func (srv *GlobalService) GetBoundContext(accountKey string) (*models.ContextBoundEvent, error) {
	subject := fmt.Sprintf("natster.events.%s.none.none.%s", accountKey, models.ContextBoundEventType)
	js, err := jetstream.New(srv.nc)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	msg, err := s.GetLastMsgForSubject(ctx, subject)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("context bound event not found")
	}

	var discoveredContext models.ContextBoundEvent
	err = json.Unmarshal(msg.Data, &discoveredContext)
	if err != nil {
		slog.Error("Deserialization failure of context bound event", err)
		return nil, err
	}

	return &discoveredContext, nil
}

func (srv *GlobalService) GetOAuthIdForAccount(accountKey string) (*string, error) {
	discoveredContext, err := srv.GetBoundContext(accountKey)
	if err != nil {
		return nil, err
	}

	return &discoveredContext.OAuthIdentity, nil
}

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
	discoveredContext := <-ch
	cc.Stop()

	return discoveredContext, nil

}

func (srv *GlobalService) GetMyCatalogs(myAccountKey string) ([]models.CatalogShareSummary, error) {
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.CatalogSharedEventType)
	js, err := jetstream.New(srv.nc)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	opts := make([]jetstream.StreamInfoOpt, 0)
	opts = append(opts, jetstream.WithSubjectFilter(subject))
	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return nil, err
	}
	streamInfo, err := s.Info(ctx, opts...)
	if err != nil {
		return nil, err
	}

	summaries := make([]models.CatalogShareSummary, 0)
	for k := range streamInfo.State.Subjects {
		tokens := strings.Split(k, ".")
		from := tokens[2]
		to := tokens[3]
		catalog := tokens[4]
		online := srv.IsCatalogOnline(catalog)

		if from == myAccountKey || to == myAccountKey {
			summaries = append(summaries, models.CatalogShareSummary{
				FromAccount:   from,
				ToAccount:     to,
				Catalog:       catalog,
				CatalogOnline: online,
			})
		}
	}

	summaries = append(summaries, models.CatalogShareSummary{
		FromAccount:   synadiaHubAccount,
		ToAccount:     myAccountKey,
		Catalog:       "synadiahub",
		CatalogOnline: true,
	})

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
	}
}
