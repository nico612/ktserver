package captcha

import (
	"context"
	"ktserver/internal/pkg/bizerr"
	"time"
)

const devModelSupperCode = "888888"

type Options struct {
	DevMode       bool
	MaxErrorCount int // 验证码最大错误次数
	Expiration    time.Duration
}

func NewOptions() *Options {
	return &Options{
		DevMode:       false,
		MaxErrorCount: 10,
		Expiration:    5 * time.Minute,
	}
}

type CaptchaService struct {
	generator Generator
	limiter   Limiter
	store     Store
	opts      *Options
}

func NewService(generator Generator, limiter Limiter, store Store, opts *Options) *CaptchaService {
	return &CaptchaService{
		generator: generator,
		limiter:   limiter,
		opts:      opts,
		store:     store,
	}
}

// Generate 生成图形验证码 id: 唯一标识，一般为ip地址，比如登录时使用ip地址限制访问频率
func (s *CaptchaService) Generate(ctx context.Context, id string) (captchaId, b64s, answer string, err error) {
	// 是否限制生成
	valid, err := s.limiter.Acquire(ctx, id)
	if err != nil {
		return "", "", "", err
	}
	// 获取太频繁
	if !valid {
		return "", "", "", bizerr.CaptchaLimit
	}

	captchaId, b64s, answer, err = s.generator.Generate()
	if err != nil {
		return "", "", "", err
	}

	code := CaptchaCode{
		Id:         captchaId,
		Value:      answer,
		Created:    time.Now(),
		Expiration: s.opts.Expiration,
	}

	if err = s.store.Set(ctx, captchaId, code); err != nil {
		return "", "", "", err
	}

	return captchaId, b64s, answer, nil
}

// Verify 验证图形验证码
func (s *CaptchaService) Verify(ctx context.Context, captchaId, answer string, clear bool) error {
	if s.opts.DevMode && answer == devModelSupperCode {
		return nil
	}

	code, err := s.store.Get(ctx, captchaId)
	if err != nil {
		return bizerr.CaptchaInvalid
	}

	err = s.check(answer, code)
	if err != nil {
		code.ErrorCount++
		// 从新保存code，记录错误次数。不刷新过期时间
		if err := s.store.Save(ctx, captchaId, code); err != nil {
			return err
		}
		return err
	}

	if clear {
		if err = s.store.Delete(ctx, captchaId); err != nil {
			return err
		}
	}

	return nil
}

// check 验证码校验
func (s *CaptchaService) check(value string, code CaptchaCode) error {

	if code.isExpired() {
		return bizerr.CaptchaExpired
	}

	// 验证码错误次数
	if code.ErrorCount >= s.opts.MaxErrorCount {
		return bizerr.CaptchaBlocked
	}

	if value != code.Value {
		return bizerr.CaptchaInvalid
	}
	return nil
}
