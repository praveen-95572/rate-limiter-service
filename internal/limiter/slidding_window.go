package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/praveen-95572/rate-limiter-service/internal/adaptor/persistance"
	"github.com/redis/go-redis/v9"
)

type SlidingWindow struct {
	store  *persistance.RedisStore
	window time.Duration
	limit  int
}

func NewSlidingWindow(store *persistance.RedisStore, window time.Duration, limit int) *SlidingWindow {
	return &SlidingWindow{
		store:  store,
		window: window,
		limit:  limit,
	}
}

func (sw *SlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()
	windowStart := now - int64(sw.window.Seconds())

	pipe := sw.store.Client.TxPipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	count := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, sw.window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return count.Val() <= int64(sw.limit), nil
}
