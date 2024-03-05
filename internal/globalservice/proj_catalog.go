package globalservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	catalogProjectionBucketName = "PROJ_CATALOG"
)

type catalogProjection struct {
	Name       string   `json:"catalog_name"`
	Owner      string   `json:"owner"`
	SharedWith []string `json:"shared_with"`
}

func (srv *GlobalService) createOrReuseCatalogProjectionConsumer() error {
	js, _ := jetstream.New(srv.nc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return err
	}
	consumer, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    "PROJ_CATALOG",
		Durable: "PROJ_CATALOG",
		FilterSubjects: []string{
			"natster.events.*.*.*.catalog_shared",
			"natster.events.*.*.*.catalog_unshared",
		},
	})
	if err != nil {
		slog.Error("Failed to create JetStream consumer for the event stream", slog.Any("error", err))
		return err
	}

	go srv.projectCatalog(consumer)
	return nil
}

func (srv *GlobalService) createOrReuseCatalogProjectionBucket() (jetstream.KeyValue, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, catalogProjectionBucketName)
	if err != nil {
		kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:       catalogProjectionBucketName,
			Description:  "Derived catalog status",
			Storage:      jetstream.FileStorage,
			MaxValueSize: maxBytes,
			MaxBytes:     maxBytes,
		})
		if err != nil {
			slog.Error("Failed to create catalog projection bucket", err)
			return nil, err
		}
		return kv, nil
	}

	return kv, nil
}

func (srv *GlobalService) projectCatalog(consumer jetstream.Consumer) {
	slog.Info("Started catalog projector consumer")

	for {
		msgs, err := consumer.Fetch(1)
		if err != nil {
			slog.Warn("Failed to fetch message from catalog projector consumer", slog.Any("error", err))
			continue
		}
		msg := <-msgs.Messages()
		srv.updateCatalogProjection(msg)
	}
}

func (srv *GlobalService) updateCatalogProjection(msg jetstream.Msg) {
	if msg == nil {
		return
	}

	tokens := strings.Split(msg.Subject(), ".")
	from := tokens[2]
	to := tokens[3]
	catalog := tokens[4]
	event_type := tokens[5]

	slog.Info("Updating catalog projection",
		slog.String("catalog", catalog),
		slog.String("from", from),
		slog.String("to", to),
		slog.String("event_type", event_type),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	js, _ := jetstream.New(srv.nc)
	kv, err := js.KeyValue(ctx, catalogProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate catalog projection bucket", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	var projection catalogProjection
	entry, err := kv.Get(ctx, catalog)
	if err != nil {
		projection = catalogProjection{
			Name:       catalog,
			Owner:      from,
			SharedWith: make([]string, 0),
		}
	} else {
		err = json.Unmarshal(entry.Value(), &projection)
		if err != nil {
			slog.Error("Corrupt projection", slog.Any("error", err))
			_ = msg.Nak()
			return
		}
	}

	switch event_type {
	case models.CatalogSharedEventType:
		if !slices.Contains(projection.SharedWith, to) {
			projection.SharedWith = append(projection.SharedWith, to)
		}
	case models.CatalogUnsharedEventType:
		idx := slices.Index(projection.SharedWith, catalog)
		if idx > -1 {
			projection.SharedWith = slices.Delete(projection.SharedWith, idx, idx+1)
		}
	}

	projBytes, err := json.Marshal(projection)
	if err != nil {
		slog.Error("Failed to serialize projection", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	_, err = kv.Put(ctx, catalog, projBytes)
	if err != nil {
		slog.Error("Failed to write projection", slog.Any("error", err))
		_ = msg.Nak()
		return
	}

	_ = msg.Ack()
}
