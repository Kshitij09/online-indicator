package redisstore

import (
	"context"
	"errors"
	"github.com/Kshitij09/online-indicator/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type lastSeenStorage struct {
	client          *redis.Client
	ctx             context.Context
	onlineThreshold time.Duration
}

func LastSeenDao(client *redis.Client, ctx context.Context, onlineThreshold time.Duration) domain.LastSeenDao {
	return lastSeenStorage{
		client:          client,
		ctx:             ctx,
		onlineThreshold: onlineThreshold,
	}
}

func (ctx lastSeenStorage) GetLastSeen(accountId string) (int64, error) {
	val, err := ctx.client.Get(ctx.ctx, accountId).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, domain.ErrSessionExpired
	}
	return val, err
}

func (ctx lastSeenStorage) SetLastSeen(accountId string, lastSeen int64) error {
	return ctx.client.Set(ctx.ctx, accountId, lastSeen, ctx.onlineThreshold).Err()
}
