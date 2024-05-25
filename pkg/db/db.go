package db

import (
	"github.com/google/wire"
)

//var ProviderSet = wire.NewSet(NewMySQL, NewRedis, wire.Bind(new(redis.UniversalClient), new(*redis.Client)))

var ProviderSet = wire.NewSet(NewMySQL, NewRedis)
