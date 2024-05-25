package service

import (
	"github.com/google/wire"
	"ktserver/internal/pkg/initdb"
	"ktserver/internal/pkg/initdb/handler"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewBaseService,
	wire.NewSet(
		handler.NewMysqlInitHandler,
		wire.Bind(new(initdb.TypedDBInitHandler), new(*handler.MysqlInitHandler)),
		initdb.NewInitDBService,
	),
	NewDBUService,
	NewUserService,
)
