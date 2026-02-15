package persistance

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisStore{Client: rdb}
}

func (r *RedisStore) Eval(
	ctx context.Context,
	script string,
	keys []string,
	args ...interface{},
) (interface{}, error) {
	cmd := r.Client.Eval(ctx, script, keys, args...)
	result, err := cmd.Result()
	if err == redis.Nil {
		return int64(0), nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisStore) Close() error {
	return r.Client.Close()
}

func (r *RedisStore) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisStore) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return r.Client.Expire(ctx, key, ttl).Err()
}
