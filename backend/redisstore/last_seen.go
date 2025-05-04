package redisstore

import (
	"context"
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
	return ctx.client.Get(ctx.ctx, accountId).Int64()
}

func (ctx lastSeenStorage) SetLastSeen(accountId string, lastSeen int64) error {
	return ctx.client.Set(ctx.ctx, accountId, lastSeen, ctx.onlineThreshold).Err()
}
