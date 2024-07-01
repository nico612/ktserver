package genericoptions

import "time"

// CaptchaOptions 图形验证码配置
type CaptchaOptions struct {
	ImageWidth    int           `json:"image-width,omitempty" mapstructure:"image-width,omitempty"`
	ImageHeight   int           `json:"image-height,omitempty" mapstructure:"image-height,omitempty"`
	KeyLength     int           `json:"key-length,omitempty" mapstructure:"key-length,omitempty"`
	PeriodSeconds int           `json:"period-seconds,omitempty" mapstructure:"period-seconds,omitempty"`
	MaxCount      int           `json:"max-count,omitempty" mapstructure:"max-count,omitempty"`
	MaxErrorCount int           `json:"max-error-count,omitempty" mapstructure:"max-error-count,omitempty"`
	BlockDuration time.Duration `json:"block-duration,omitempty" mapstructure:"block-duration,omitempty"`
	Expiration    time.Duration `json:"expiration,omitempty" mapstructure:"expiration,omitempty"`
}

// NewCaptchaOptions create a `zero` value instance.
func NewCaptchaOptions() *CaptchaOptions {
	return &CaptchaOptions{
		ImageWidth:    240,
		ImageHeight:   80,
		KeyLength:     5,
		PeriodSeconds: 60,               // 60 seconds
		MaxCount:      3,                // 100 times 60秒内最多生成10次验证码
		MaxErrorCount: 5,                // 验证码最大错误次数
		BlockDuration: 60 * time.Minute, // 锁定时长 5分钟
		Expiration:    5 * time.Minute,  // 10分钟过期
	}
}
