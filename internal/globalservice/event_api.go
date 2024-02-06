package globalservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Returns the total number of accounts in which developers have run `natster init`
func (srv *GlobalService) GetTotalInitializedAccounts() (uint64, error) {
	js, err := jetstream.New(srv.nc)
	if err != nil {
		return 0, err
	}
	ctx := context.Background()
	subject := fmt.Sprintf("natster.events.*.*.*.%s", models.NatsterInitializedEventType)
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
	return streamInfo.State.Msgs, nil
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
		err = srv.nc.Publish(subject, raw)
		if err != nil {
			slog.Error("Failed to publish Natster event", err)
			return
		}
	}
}
