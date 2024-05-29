package admin

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jinzhu/copier"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/data"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/genericoptions"
	"ktserver/internal/pkg/source"
	"ktserver/pkg/db"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "admin-server"
	// Version is the version of the compiled software.
	Version string
)

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
	)
}

func NewApp() {

	c := conf.ViperParse(Name)

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.name", Name,
		"service.version", Version,
		//"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
	)

	var mysqlOptions db.MySQLOptions
	var redisOptions db.RedisOptions

	_ = copier.Copy(&mysqlOptions, &c.MySQLOptions)
	_ = copier.Copy(&redisOptions, &c.RedisOptions)

	// jwt options
	var jwtOpts auth.JWTOpts
	var authOpts genericoptions.AuthOptions
	_ = copier.Copy(&jwtOpts, &c.AuthOptions)
	_ = copier.Copy(&authOpts, &c.AuthOptions)

	// initialize the application
	app, cleanup, err := wireApp(c, &mysqlOptions, &redisOptions, &jwtOpts, &authOpts, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 数据库迁移
	if err := data.AutoMigrate(); err != nil {
		logger.Log(log.LevelError, "AutoMigrate", err)
	}

	// 注册需要初始化的系统数据
	source.RegisterSource()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
