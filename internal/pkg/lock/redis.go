package lock

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"ktserver/internal/pkg/bizerr"
	"math/rand"
	"strings"
	"time"
)

type RedisLocker struct {
	log *log.Helper
	rds redis.UniversalClient
}

func NewRedisLocker(logger log.Logger, rds redis.UniversalClient) *RedisLocker {
	return &RedisLocker{
		log: log.NewHelper(logger),
		rds: rds,
	}
}

func (l *RedisLocker) WithLock(ctx context.Context, key string, holdingTime time.Duration) (unlock func(), err error) {
	if len(key) > 100 {
		return nil, bizerr.InvalidParam.WithMsg("key too long")
	}

	locker := newRedisLock(l.log, l.rds, key, holdingTime)
	if err := locker.Lock(ctx); err != nil {
		return nil, err
	}

	unlock = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		unlocked, err := locker.Unlock(ctx)
		if err != nil {

			l.log.Error("unlock error")
		}
		if !unlocked {
			l.log.Warnw("key", locker.key, "id", locker.id, "error", "invalid unlock")
		}
	}
	return
}

type redisLock struct {
	log         *log.Helper
	rds         redis.UniversalClient
	key         string
	id          string
	holdingTime time.Duration
}

const redisLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`

const redisUnlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`

var redisLockRng = rand.New(rand.NewSource(time.Now().UnixNano()))

func newRedisLock(logger *log.Helper, rds redis.UniversalClient, key string, holdingTime time.Duration) *redisLock {
	const alpha = "abcdefghijklmnopqrstuvwzyxABCDEFGHIJKLMNOPQRSTUVWZYX0123456789"
	var id strings.Builder
	for i := 0; i < 16; i++ {
		id.WriteRune(rune(alpha[redisLockRng.Intn(len(alpha))]))
	}

	return &redisLock{
		log:         logger,
		rds:         rds,
		id:          id.String(),
		key:         "lock:" + key,
		holdingTime: holdingTime,
	}
}

// Lock 获取锁，直到 context 超时
func (l *redisLock) Lock(ctx context.Context) error {
	tries := 0

	for {
		select {
		case <-ctx.Done():
			return errors.WithStack(ctx.Err())
		default:
			locked, err := l.TryLock(ctx)
			if err != nil {
				return err
			}

			if !locked {
				// 没有获取到锁，自旋重试
				tries++
				sleepTime := 0.01 * float64(tries*tries) * 1000 // 0.01*2^x ms 10,40,90,160,250,360,490,640 ...
				if sleepTime > 2000 {
					sleepTime = 2000
				}

				waiting := time.NewTimer(time.Duration(sleepTime) * time.Millisecond)
				select {
				case <-ctx.Done():
					waiting.Stop()
					return errors.WithStack(ctx.Err())
				case <-waiting.C:
				}

				continue
			}

			// 获取到锁
			return nil
		}
	}
}

// TryLock 获取锁
func (l *redisLock) TryLock(ctx context.Context) (bool, error) {
	cmd := l.rds.Eval(ctx, redisLockScript, []string{l.key}, []string{
		l.id, cast.ToString(l.holdingTime.Milliseconds()),
	})

	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, errors.Wrap(err, "lock error")
	}

	if v, ok := cmd.Val().(string); ok && v == "OK" {
		return true, nil
	}

	return false, bizerr.UnknownError.WithMsgf("lock error:%v", cmd.Val())
}

func (l *redisLock) Unlock(ctx context.Context) (bool, error) {
	cmd := l.rds.Eval(ctx, redisUnlockScript, []string{l.key}, []string{l.id})
	if cmd.Err() != nil {
		return false, errors.Wrap(cmd.Err(), "unlock error")
	}

	if v, ok := cmd.Val().(int64); ok && v == 1 {
		return true, nil
	}

	return false, nil
}
