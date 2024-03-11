package globalservice

import (
	"context"

	"log/slog"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/synadia-labs/natster/internal/models"
)

type GlobalService struct {
	nc      *nats.Conn
	hbCache *ttlcache.Cache[string, models.Heartbeat]
}

func New(nc *nats.Conn) *GlobalService {
	return &GlobalService{
		nc: nc,
		hbCache: ttlcache.New[string, models.Heartbeat](
			ttlcache.WithTTL[string, models.Heartbeat](3 * time.Minute)),
	}
}

func (srv *GlobalService) Start(version, commit, date string) error {
	err := srv.startApiSubscriptions()
	if err != nil {
		return err
	}

	err = srv.createOrReuseEventStream()
	if err != nil {
		return err
	}

	_, err = srv.createOrReuseOtcBucket()
	if err != nil {
		return err
	}

	_, err = srv.createOrReuseCatalogProjectionBucket()
	if err != nil {
		return err
	}

	_, err = srv.createOrReuseAccountProjectionBucket()
	if err != nil {
		return err
	}

	err = srv.createOrReuseCatalogProjectionConsumer()
	if err != nil {
		return err
	}

	err = srv.createOrReuseAccountProjectionConsumer()
	if err != nil {
		return err
	}

	srv.hbCache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, models.Heartbeat]) {
		if reason == ttlcache.EvictionReasonCapacityReached {
			slog.Info("Evicting heartbeat cache item - capacity reached",
				slog.String("key", item.Key()))
		} else if reason == ttlcache.EvictionReasonDeleted {
			slog.Info("Evicting heartbeat cache item - deleted",
				slog.String("key", item.Key()))
		} else if reason == ttlcache.EvictionReasonExpired {
			slog.Info("Evicting heartbeat cache item - expired",
				slog.String("key", item.Key()))
		}
	})

	go srv.hbCache.Start()

	slog.Info("Natster Global Service Started",
		slog.String("version", version),
		slog.String("commit", commit),
		slog.String("date", date),
	)
	return nil
}

func (srv *GlobalService) Stop() error {
	srv.nc.Drain()
	return nil
}

func (srv *GlobalService) CreateKeyValueContext() (jetstream.KeyValue, error) {
	js, _ := jetstream.New(srv.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	kv, err := js.KeyValue(ctx, accountProjectionBucketName)
	if err != nil {
		slog.Error("Failed to locate key value bucket", slog.Any("error", err), slog.String("bucket", accountProjectionBucketName))
		return nil, err
	}
	return kv, nil
}
