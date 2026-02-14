package benchmarks

import (
	"context"
	"testing"

	"github.com/praveen-95572/rate-limiter-service/internal/adaptor/persistance"
	"github.com/praveen-95572/rate-limiter-service/internal/limiter"
)

func BenchmarkTokenBucket(b *testing.B) {
	store := persistance.NewRedisStore("localhost:6379")
	tb := limiter.NewTokenBucket(store, 1000, 100)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow(ctx, "bench_user")
	}
}
