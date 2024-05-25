package limiter

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"net/netip"
	"time"
)

// Limiter limiter
type Limiter interface {
	// Acquire id 标识, periodSeconds 限制周期, maxCount 限制次数， 如果在周期内超过限制次数，返回 false
	Acquire(ctx context.Context, id string, periodSeconds, maxCount int) (bool, error)
}

// Blocker blocker
type Blocker interface {
	IsBlocked(ctx context.Context, id string) (bool, error)          // 是否被限制
	Block(ctx context.Context, id string, until time.Duration) error // 限制
}

type Rule struct {
	Prefix        string
	PeriodSeconds int // 限制周期,秒
	MaxCount      int // 限制次数
}

func ip2long(ip string) (uint32, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if !addr.Is4() {
		return 0, errors.New(fmt.Sprintf("not ip v4: %v", addr))
	}

	b, err := addr.MarshalBinary()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return binary.BigEndian.Uint32(b), nil
}
