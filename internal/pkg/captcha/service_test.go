package captcha

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/magiconair/properties/assert"
	"github.com/mojocn/base64Captcha"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/genericoptions"
	"ktserver/internal/pkg/limiter"
	"ktserver/pkg/db"
	"testing"
)

func TestBase64Generator(t *testing.T) {
	driver := base64Captcha.DefaultDriverDigit
	generator := NewBase64Generator(driver)

	id, b64s, answer, err := generator.Generate()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("id: %s, b64s: %s, answer: %s\n", id, b64s, answer)
}

func getCaptchaService(caOpts *genericoptions.CaptchaOptions) (*CaptchaService, error) {
	opts := genericoptions.NewRedisOptions()
	var rdsO db.RedisOptions
	_ = copier.Copy(&rdsO, opts)
	rds, err := db.NewRedis(&rdsO)
	if err != nil {
		return nil, err
	}
	// 验证码生成器
	driver := base64Captcha.NewDriverDigit(caOpts.ImageHeight, caOpts.ImageWidth, caOpts.KeyLength, 0.7, 80)
	generator := NewBase64Generator(driver)

	// 放爆器
	guard := NewGuard(
		limiter.NewRedisLimiter(rds),
		limiter.NewRedisBlocker(rds),
		[]limiter.Rule{ // 多少秒内最多生成多少次验证码
			{
				PeriodSeconds: caOpts.PeriodSeconds,
				MaxCount:      caOpts.MaxCount,
			},
		},
		caOpts.BlockDuration, // 锁定时长
	)
	store := NewRedisStore(rds)
	captchaService := NewService(generator, guard, store, &Options{
		DevMode:       false,
		MaxErrorCount: caOpts.MaxErrorCount,
		Expiration:    caOpts.Expiration,
	})

	return captchaService, nil
}

func TestCaptchaService_Generate(t *testing.T) {

	caOpts := genericoptions.NewCaptchaOptions()

	captchaService, err := getCaptchaService(caOpts)
	if err != nil {
		t.Errorf("getCaptchaService error: %v", err)
	}
	ctx := context.Background()

	id, _, answer, err := captchaService.Generate(ctx, "2222")
	if err != nil {
		t.Errorf("Generate error: %v", err)
	}

	if err = captchaService.Verify(ctx, id, answer, true); err != nil {
		t.Errorf("Verify error: %v", err)
	}
}

func TestCaptchaService_Limiter(t *testing.T) {

	caOpts := genericoptions.NewCaptchaOptions()

	captchaService, err := getCaptchaService(caOpts)
	if err != nil {
		t.Errorf("getCaptchaService error: %v", err)
	}

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		_, _, _, err := captchaService.Generate(ctx, "11111")
		if i > caOpts.MaxCount {
			assert.Equal(t, err, bizerr.CaptchaLimit)
		}
	}
}

func TestCaptchaService_Block(t *testing.T) {
	caopts := genericoptions.NewCaptchaOptions()
	captchaService, err := getCaptchaService(caopts)
	if err != nil {
		t.Errorf("getCaptchaService error: %v", err)
	}

	ctx := context.Background()
	id, _, _, err := captchaService.Generate(ctx, "66666")
	if err != nil {
		t.Errorf("Generate error: %v", err)
	}

	var verr error
	// 验证错误
	for i := 0; i < 10; i++ {
		verr = captchaService.Verify(ctx, id, "wrong", true)
		if i >= caopts.MaxErrorCount {
			assert.Equal(t, verr, bizerr.CaptchaBlocked)
		} else {
			assert.Equal(t, verr, bizerr.CaptchaInvalid)
		}
	}
}
