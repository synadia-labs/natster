package globalservice

import (
	"context"

	"log/slog"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/nats-io/nats.go"
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

func (srv *GlobalService) Start() error {
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
	slog.Info("Natster Global Service Started")
	return nil
}

func (srv *GlobalService) Stop() error {
	srv.nc.Drain()
	return nil
}
