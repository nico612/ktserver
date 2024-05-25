package conf

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/spf13/viper"
	"ktserver/internal/pkg/genericoptions"
)

type Config struct {
	Server         *genericoptions.ServerOptions  `mapstructure:"server" json:"server" yaml:"server"`
	MySQLOptions   *genericoptions.MySQLOptions   `json:"mysql" mapstructure:"mysql"`
	RedisOptions   *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	CaptchaOptions *genericoptions.CaptchaOptions `json:"captcha" mapstructure:"captcha"`
	AuthOptions    *genericoptions.AuthOptions    `json:"auth" mapstructure:"auth"`
}

func NewConfig() *Config {
	opts := &Config{
		Server:         genericoptions.NewServerOptions(),
		MySQLOptions:   genericoptions.NewMySQLOptions(),
		RedisOptions:   genericoptions.NewRedisOptions(),
		CaptchaOptions: genericoptions.NewCaptchaOptions(),
		AuthOptions:    genericoptions.NewAuthOptions(),
	}
	return opts
}

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// Parse parse config file
func Parse(appName string) *Config {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	opts := NewConfig()
	if err := c.Scan(opts); err != nil {
		panic(err)
	}
	return opts
}

func ViperParse(appName string) *Config {
	flag.Parse()

	viper.AutomaticEnv()
	viper.SetConfigFile(flagconf)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var c Config
	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变化时的回调
		// 重新解析配置文件
		if err := viper.Unmarshal(&c); err != nil {
			panic(err)
		}
	})

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return &c
}
