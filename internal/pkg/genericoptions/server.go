package genericoptions

import "time"

type Mode string

const (
	ModeDev  Mode = "dev"  // 开发模式
	ModeTest Mode = "test" // 测试模式
	ModeProd Mode = "prod" // 生产模式
)

// IsDevMode 是否是开发模式
func (m Mode) IsDevMode() bool {
	return m == ModeDev
}

type ServerOptions struct {
	Mode  Mode         `mapstructure:"mode" json:"mode" yaml:"mode"` // 模式
	Admin *HttpOptions `mapstructure:"admin" json:"admin" yaml:"admin"`
	App   *HttpOptions `mapstructure:"app" json:"app" yaml:"app"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Mode:  ModeDev,
		Admin: NewHttpOptions(),
		App:   NewHttpOptions(),
	}
}

type HttpOptions struct {
	Network string        `mapstructure:"network" json:"network" yaml:"network"` // 网络
	Addr    string        `mapstructure:"addr" json:"addr" yaml:"addr"`          // 地址
	Timeout time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"` // 超时
}

func NewHttpOptions() *HttpOptions {
	return &HttpOptions{
		Network: "tcp",
		Addr:    "127.0.0.1:8000",
		Timeout: 10,
	}
}
