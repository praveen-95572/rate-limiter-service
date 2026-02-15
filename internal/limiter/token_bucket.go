package limiter

import (
	"context"
	"log"
	"time"

	"github.com/praveen-95572/rate-limiter-service/internal/adaptor/persistance"
)

type TokenBucket struct {
	store    *persistance.RedisStore
	capacity int
	rate     int // tokens per second
}

func NewTokenBucket(store *persistance.RedisStore, capacity, rate int) *TokenBucket {
	return &TokenBucket{
		store:    store,
		capacity: capacity,
		rate:     rate,
	}
}

var tokenBucketScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "timestamp")
local tokens = tonumber(data[1])
local timestamp = tonumber(data[2])

if tokens == nil or timestamp == nil then
	tokens = capacity
	timestamp = now
end

local delta = math.max(0, now - timestamp)
local refill = delta * rate
tokens = math.min(capacity, tokens + refill)
local allowed = tokens >= requested

if allowed then
	tokens = tokens - requested
end

redis.call("HMSET", key, "tokens", tokens, "timestamp", now)
redis.call("EXPIRE", key, 3600)

return allowed
`

func (tb *TokenBucket) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()

	result, err := tb.store.Eval(
		ctx,
		tokenBucketScript,
		[]string{key},
		tb.capacity,
		tb.rate,
		now,
		1,
	)

	if err != nil {
		log.Printf("Error for key %s is: %v", key, err)
		return false, err
	}

	return result.(int64) == 1, nil
}
