package captcha

import (
	"context"
	"ktserver/internal/pkg/limiter"
	"time"
)

type Guard struct {
	blocker       limiter.Blocker
	limiter       limiter.Limiter
	rules         []limiter.Rule // 限制规则，这里设计为多个规则，方便处理多种情况
	blockDuration time.Duration  // 限制时间
}

func NewGuard(
	limiter limiter.Limiter,
	blocker limiter.Blocker,
	rules []limiter.Rule,
	blockDuration time.Duration,
) *Guard {
	return &Guard{
		blocker:       blocker,
		limiter:       limiter,
		rules:         rules,
		blockDuration: blockDuration,
	}
}

func (l *Guard) Acquire(ctx context.Context, id string) (bool, error) {
	key := storeKey(id)
	blocked, err := l.blocker.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}

	if blocked {
		return false, nil
	}

	for _, r := range l.rules {

		ac, err := l.limiter.Acquire(ctx, key, r.PeriodSeconds, r.MaxCount)
		if err != nil {
			return false, err
		}
		if !ac {
			err := l.blocker.Block(ctx, key, l.blockDuration)
			if err != nil {
				return false, err
			}

			return false, nil
		}
	}

	return true, nil

}
