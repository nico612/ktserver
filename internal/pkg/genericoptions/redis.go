package genericoptions

import "time"

type RedisOptions struct {
	Addr         string        `json:"addr" mapstructure:"addr"`
	Username     string        `json:"username" mapstructure:"username"`
	Password     string        `json:"password" mapstructure:"password"`
	Database     int           `json:"database" mapstructure:"database"`
	MaxRetries   int           `json:"max-retries" mapstructure:"max-retries"`
	MinIdleConns int           `json:"min-idle-conns" mapstructure:"min-idle-conns"`
	DialTimeout  time.Duration `json:"dial-timeout" mapstructure:"dial-timeout"`
	ReadTimeout  time.Duration `json:"read-timeout" mapstructure:"read-timeout"`
	WriteTimeout time.Duration `json:"write-timeout" mapstructure:"write-timeout"`
	PoolTimeout  time.Duration `json:"pool-time" mapstructure:"pool-time"`
	PoolSize     int           `json:"pool-size" mapstructure:"pool-size"`
	// tracing switch
	EnableTrace bool `json:"enable-trace" mapstructure:"enable-trace"`
}

// NewRedisOptions create a `zero` value instance.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addr:         "127.0.0.1:6379",
		Username:     "",
		Password:     "",
		Database:     0,
		MaxRetries:   3,
		MinIdleConns: 0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		EnableTrace:  false,
	}
}
