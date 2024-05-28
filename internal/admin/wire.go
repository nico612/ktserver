//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package admin

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"ktserver/internal/admin/biz"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/data"
	"ktserver/internal/admin/server"
	"ktserver/internal/admin/service"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/authz"
	"ktserver/internal/pkg/genericoptions"
	"ktserver/internal/pkg/lock"
	"ktserver/internal/pkg/middleware"
	"ktserver/pkg/db"
)

// wireApp init kratos application.
func wireApp(*conf.Config, *db.MySQLOptions, *db.RedisOptions, *auth.JWTOpts, *genericoptions.AuthOptions, log.Logger) (*kratos.App, func(), error) {

	panic(
		wire.Build(
			wire.NewSet(
				lock.NewRedisLocker,
				auth.NewDefaultTokener,
				auth.NewAuthenticator,
				authz.NewCasbinAuthorizer,
				middleware.NewJWTMiddleware,
				middleware.NewAuthzMiddleware,
			),
			db.ProviderSet,
			biz.ProviderSet,
			data.ProviderSet,
			service.ProviderSet,
			server.ProviderSet,
			newApp,
		),
	)
}
