package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"ktserver/internal/admin/conf"
	"ktserver/internal/pkg/captcha"
	"ktserver/internal/pkg/limiter"
)

type CaptchaUseCase struct {
	c       *conf.Config
	log     *log.Helper
	rds     redis.UniversalClient
	captcha *captcha.CaptchaService
}

func NewCaptchaUseCase(c *conf.Config, logger log.Logger, rds redis.UniversalClient) *CaptchaUseCase {
	// 验证码生成器
	driver := base64Captcha.NewDriverDigit(c.CaptchaOptions.ImageHeight, c.CaptchaOptions.ImageWidth, c.CaptchaOptions.KeyLength, 0.7, 80)
	generator := captcha.NewBase64Generator(driver)

	// 防爆器
	guard := captcha.NewGuard(
		limiter.NewRedisLimiter(rds),
		limiter.NewRedisBlocker(rds),
		[]limiter.Rule{ // 多少秒内最多生成多少次验证码
			{
				PeriodSeconds: c.CaptchaOptions.PeriodSeconds,
				MaxCount:      c.CaptchaOptions.MaxCount,
			},
		},
		c.CaptchaOptions.BlockDuration, // 锁定时长
	)

	opts := &captcha.Options{
		DevMode:       c.Server.Mode.IsDevMode(),
		MaxErrorCount: c.CaptchaOptions.MaxErrorCount,
		Expiration:    c.CaptchaOptions.Expiration,
	}

	captchaStore := captcha.NewRedisStore(rds)
	captcha := captcha.NewService(generator, guard, captchaStore, opts)

	return &CaptchaUseCase{
		c:       c,
		log:     log.NewHelper(logger),
		captcha: captcha,
	}
}

// GenerateCaptcha 生成验证码
func (uc *CaptchaUseCase) GenerateCaptcha(ctx context.Context, ipAddr string) (captchaId, b64s, answer string, err error) {
	return uc.captcha.Generate(ctx, ipAddr)
}

// VerifyCaptcha 验证验证码
func (uc *CaptchaUseCase) VerifyCaptcha(ctx context.Context, captchaId, answer string) error {
	return uc.captcha.Verify(ctx, captchaId, answer, true)
}
