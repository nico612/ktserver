package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/service"
	"ktserver/internal/pkg/middleware"
)

type Router struct {
	g      *gin.Engine
	c      *conf.Config
	prefix string
	base   *service.BaseService
	initdb *service.DBUService
	user   *service.UserService
	authmw *middleware.JWTMiddleware
}

func NewRouter(
	c *conf.Config,
	base *service.BaseService,
	initdb *service.DBUService,
	user *service.UserService,
	authmw *middleware.JWTMiddleware,
) *Router {
	r := &Router{
		g:      gin.Default(),
		c:      c,
		prefix: "/api",
		base:   base,
		initdb: initdb,
		user:   user,
		authmw: authmw,
	}
	return r
}

func (r *Router) installRouter() {
	g := r.g
	g.Use(gin.Recovery())

	//if r.c.Server.Mode != genericoptions.ModeProd {
	//	g.Use(gin.Logger())
	//}

	g.GET(r.prefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//publicGroup := g.Group("")
	publicGroup := g.Group(r.prefix)
	{
		// 健康监测
		publicGroup.GET("health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}

	// 基础路由，提供登录注册
	baseRouter := publicGroup.Group("base")
	{
		baseRouter.POST("login", r.base.Login)     // 登录
		baseRouter.POST("captcha", r.base.Captcha) // 生成验证码
	}

	// 初始化数据库
	initDBRouter := publicGroup.Group("init")
	{
		initDBRouter.POST("initdb", r.initdb.InitDB)   // 初始化数据库
		initDBRouter.POST("checkdb", r.initdb.CheckDB) // 检测是否需要初始化数据库
	}

	// 私有路由，需要登录
	privateGroup := g.Group(r.prefix)
	userRouter := privateGroup.Group("user").Use(r.authmw.AuthFunc())
	{
		userRouter.GET("getUserInfo", r.user.GetUserInfo)
		userRouter.POST("changePassword", r.user.ChangePassword)
	}

	//privateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//apiRouter := privateGroup.Group("api")

}
