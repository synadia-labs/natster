package globalservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName    = "NATSTER_EVENTS"
	streamSubject = "natster.events.*.*.*.*"
	maxBytes      = 1_073_741_824 // 1 gib

	otcBucketName = "NATSTER_CODES"
	otcTimeoutMinutes
)

// Stream pattern
// natster.events.{origin}.{target}.{catalog}.{event}
// e.g.
// natster.events.Axxx.Axxxxy.kevbuzz.catalog_shared
// When no target is relevant, token `none` is used
// When no catalog is relevant, token `none` is used
// e.g.
// natster.events.Axxxx.none.none.natster_initialized

func (srv *GlobalService) createOrReuseEventStream() error {
	js, _ := jetstream.New(srv.nc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := js.Stream(ctx, streamName)
	if err != nil {
		_, err = createStream(srv.nc, js)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *GlobalService) createOrReuseOtcBucket() (jetstream.KeyValue, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, otcBucketName)
	if err != nil {
		kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:       otcBucketName,
			Description:  "Natster Global Auto-Expiring OTCs",
			TTL:          tokenValidTimeMinutes * time.Minute,
			Storage:      jetstream.FileStorage,
			MaxValueSize: maxBytes,
			MaxBytes:     maxBytes,
		})
		if err != nil {
			slog.Error("Failed to create OTC bucket", err)
			return nil, err
		}
		return kv, nil
	}

	return kv, nil
}

func createStream(nc *nats.Conn, js jetstream.JetStream) (jetstream.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:        streamName,
		Description: "Natster Global Event Stream",
		Subjects:    []string{streamSubject},
		MaxBytes:    maxBytes,
	})
	if err != nil {
		return nil, err
	}
	return s, nil
}
