package globalservice

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	accountProjectionBucketName = "PROJ_ACCOUNT"
)

// Represents the current state of a single Natster account
type accountProjection struct {
	// List of catalogs this account has shared with others
	OutShares []shareEntry `json:"out_shares"`
	// List of catalogs that have been shared with this account
	InShares []shareEntry `json:"in_shares"`
	// UTC unix timestamp when the account was initialized
	InitializedAt int64                     `json:"initialized_at"`
	BoundContext  *models.ContextBoundEvent `json:"bound_context,omitempty"`
}

type shareEntry struct {
	Catalog string `json:"catalog"`
	Account string `json:"account"`
}

func (srv *GlobalService) createOrReuseAccountProjectionConsumer() error {
	js, _ := jetstream.New(srv.nc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s, err := js.Stream(ctx, streamName)
	if err != nil {
		return err
	}
	consumer, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          accountProjectionBucketName,
		Durable:       accountProjectionBucketName,
		FilterSubject: "natster.events.*.*.*.*",
	})
	if err != nil {
		slog.Error("Failed to create JetStream consumer for the event stream", slog.Any("error", err))
		return err
	}

	go srv.projectAccount(consumer)
	return nil
}

func (srv *GlobalService) createOrReuseAccountProjectionBucket() (jetstream.KeyValue, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		// Default is to create keys with only 1 value (e.g. no history)
		kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:       accountProjectionBucketName,
			Description:  "Derived account status",
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

func (srv *GlobalService) projectAccount(consumer jetstream.Consumer) {
	slog.Info("Started account projector consumer")

	for {
		msgs, err := consumer.Fetch(1)
		if err != nil {
			slog.Warn("Failed to fetch message from account projector consumer", slog.Any("error", err))
			continue
		}
		for msg := range msgs.Messages() {
			srv.updateAccountProjection(msg)
		}
	}
}

func (srv *GlobalService) updateAccountProjection(msg jetstream.Msg) {
	if msg == nil {
		return
	}

	tokens := strings.Split(msg.Subject(), ".")
	from := tokens[2]
	to := tokens[3]
	catalog := tokens[4]
	event_type := tokens[5]

	slog.Info("Updating account projection",
		slog.String("catalog", catalog),
		slog.String("from", from),
		slog.String("to", to),
		slog.String("event_type", event_type),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	js, _ := jetstream.New(srv.nc)
	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate account projection bucket", slog.Any("error", err))
		_ = msg.Nak()
		return
	}

	switch event_type {
	case models.CatalogSharedEventType:
		addShare(msg, from, to, catalog, kv)
	case models.CatalogUnsharedEventType:
		removeShare(msg, from, to, catalog, kv)
	case models.NatsterInitializedEventType:
		initAccount(msg, from, kv)
	case models.ContextBoundEventType:
		recordContextBinding(msg, from, kv)
	}
}

func recordContextBinding(msg jetstream.Msg, account string, kv jetstream.KeyValue) {
	existingAccount, err := loadAccount(kv, account)
	if err != nil {
		slog.Error("Failed to load account corresponding to context binding event", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	if existingAccount == nil {
		slog.Error("Attempted to bind context to an account that hasn't been initialized")
		_ = msg.Nak()
		return
	}
	var bindingEvent models.ContextBoundEvent
	err = json.Unmarshal(msg.Data(), &bindingEvent)
	if err != nil {
		slog.Error("Corrupt context bound event", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	existingAccount.BoundContext = &bindingEvent
	err = writeAccount(kv, account, *existingAccount)
	if err != nil {
		slog.Error("Failed to write account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}

	_ = msg.Ack()
}

func initAccount(msg jetstream.Msg, account string, kv jetstream.KeyValue) {
	newAccount, _ := loadAccount(kv, account)
	if newAccount == nil {
		newAccount = &accountProjection{
			OutShares: []shareEntry{},
			InShares:  []shareEntry{},
		}
	}
	newAccount.InitializedAt = time.Now().UTC().Unix()

	err := writeAccount(kv, account, *newAccount)
	if err != nil {
		slog.Error("Failed to write new account projection", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	_ = msg.Ack()

}

func addShare(msg jetstream.Msg, from string, to string, catalog string, kv jetstream.KeyValue) {

	if from == to {
		_ = msg.Ack()
		return
	}

	fromAccount, err := loadAccount(kv, from)
	if fromAccount == nil {
		slog.Error("Failed to load source account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	toAccount, err := loadAccount(kv, to)
	if toAccount == nil {
		slog.Error("Failed to load target account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	if !slices.ContainsFunc(fromAccount.OutShares, hasAccount(from, catalog)) {
		fromAccount.OutShares = append(fromAccount.OutShares, shareEntry{
			Account: to,
			Catalog: catalog,
		})
	}
	if !slices.ContainsFunc(toAccount.InShares, hasAccount(to, catalog)) {
		toAccount.InShares = append(toAccount.InShares, shareEntry{
			Account: from,
			Catalog: catalog,
		})
	}
	err = writeAccount(kv, from, *fromAccount)
	if err != nil {
		slog.Error("Failed to write source account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	err = writeAccount(kv, to, *toAccount)
	if err != nil {
		slog.Error("Failed to write target account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	_ = msg.Ack()
}

func hasAccount(target string, catalog string) func(s shareEntry) bool {
	return func(s shareEntry) bool {
		return s.Account == target && s.Catalog == catalog
	}
}

func removeShare(msg jetstream.Msg, from string, to string, catalog string, kv jetstream.KeyValue) {
	if from == to {
		_ = msg.Ack()
		return
	}

	fromAccount, err := loadAccount(kv, from)
	if fromAccount == nil {
		slog.Error("Failed to load source account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	toAccount, err := loadAccount(kv, to)
	if toAccount == nil {
		slog.Error("Failed to load target account", slog.Any("error", err))
		_ = msg.Nak()
		return
	}
	fromAccount.OutShares = slices.DeleteFunc(fromAccount.OutShares, hasAccount(to, catalog))
	toAccount.InShares = slices.DeleteFunc(toAccount.InShares, hasAccount(from, catalog))
	err = writeAccount(kv, from, *fromAccount)
	if err != nil {
		_ = msg.Nak()
		slog.Error("Failed to write source account", slog.Any("error", err))
		return
	}
	err = writeAccount(kv, to, *toAccount)
	if err != nil {
		_ = msg.Nak()
		slog.Error("Failed to write target account", slog.Any("error", err))
		return
	}
	_ = msg.Ack()
}

// Retrieves the account projection corresponding to the key. If the projection/key
// does not exist, (nil, nil) will be returned as this does not indicate an error
func loadAccount(kv jetstream.KeyValue, key string) (*accountProjection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var projection accountProjection
	entry, err := kv.Get(ctx, key)

	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	} else {
		err = json.Unmarshal(entry.Value(), &projection)
		if err != nil {
			slog.Error("Corrupt projection", slog.Any("error", err))
			return nil, err
		}
	}

	return &projection, nil
}

func writeAccount(kv jetstream.KeyValue, key string, projection accountProjection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	projBytes, err := json.Marshal(projection)
	if err != nil {
		slog.Error("Failed to serialize projection", slog.Any("error", err))
		return err
	}
	_, err = kv.Put(ctx, key, projBytes)
	if err != nil {
		slog.Error("Failed to write projection", slog.Any("error", err))
		return err
	}
	return nil
}
