package captcha

import (
	"context"
	"fmt"
	"time"
)

// Generator 验证码生成器接口
type Generator interface {
	Generate() (id, b64s, answer string, err error)
}

// Limiter 验证码生成限制器, id 为标识, 比如 登录时可以使用ip地址作为id
type Limiter interface {
	// Acquire 获取验证码生成权限, 返回是否可以生成验证码， 如果太频繁返回false
	Acquire(ctx context.Context, id string) (bool, error)
}

// Store 验证码存储器
type Store interface {
	Set(ctx context.Context, key string, code CaptchaCode) error
	Get(ctx context.Context, key string) (CaptchaCode, error)
	Save(ctx context.Context, key string, code CaptchaCode) error
	Delete(ctx context.Context, key string) error
}

type CaptchaCode struct {
	Id         string        `json:"id"`
	Value      string        `json:"value"`
	Created    time.Time     `json:"created"`
	Expiration time.Duration `json:"expiration"`
	ErrorCount int           `json:"error_count"`
}

// isExpired 判断验证码是否过期
func (c *CaptchaCode) isExpired() bool {
	return time.Now().After(c.Created.Add(c.Expiration))
}

func storeKey(id string) string {
	return fmt.Sprintf("captcha:%s", id)
}
