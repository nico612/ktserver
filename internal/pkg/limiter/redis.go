package limiter

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisLimiter Redis 限制器
type RedisLimiter struct {
	rds redis.UniversalClient
}

func NewRedisLimiter(rds redis.UniversalClient) *RedisLimiter {
	return &RedisLimiter{rds: rds}
}

// Acquire 获取是否应该被限制
func (r *RedisLimiter) Acquire(ctx context.Context, id string, periodSeconds, maxCount int) (bool, error) {
	nsec := time.Now().UnixMilli()

	sc := redis.NewScript(`
		local key = KEYS[1]
		local nsec = tonumber(ARGV[1])
		local periodSeconds = tonumber(ARGV[2])
		local max = tonumber(ARGV[3])

		redis.call("ZREMRANGEBYSCORE", key, "0", nsec - periodSeconds*1000)

		local count = redis.call("ZCARD", key)
		if count >= max then
			return 0
		end

		local member = nsec .. ":" .. count 

		redis.call("ZADD", key, nsec, member)
		redis.call("EXPIRE", key, periodSeconds)

		return 1
	`)

	v, err := sc.Run(ctx, r.rds, []string{r.limitKey(id)}, nsec, periodSeconds, maxCount).Int()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return v == 1, nil
}

func (r *RedisLimiter) limitKey(id string) string {
	return fmt.Sprintf("%s:%s", "limit", id)
}

// RedisBlocker Redis 阻塞器
type RedisBlocker struct {
	rds redis.UniversalClient
}

func NewRedisBlocker(rds redis.UniversalClient) *RedisBlocker {
	return &RedisBlocker{rds: rds}
}

func (r *RedisBlocker) IsBlocked(ctx context.Context, id string) (bool, error) {
	n, err := r.rds.Exists(ctx, r.blockKey(id)).Result()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return n > 0, nil
}

func (r *RedisBlocker) Block(ctx context.Context, id string, until time.Duration) error {
	_, err := r.rds.Set(ctx, r.blockKey(id), "1", until).Result()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *RedisBlocker) blockKey(id string) string {
	return fmt.Sprintf("%s:%s", "block", id)
}
