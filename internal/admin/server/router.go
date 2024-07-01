package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/service"
	"ktserver/internal/pkg/middleware"
)

type RouterService struct {
	Authmw    *middleware.JWTMiddleware   // JWT 中间件
	Authzmw   *middleware.AuthzMiddleware // 授权中间件
	Base      *service.BaseService
	Initdb    *service.DBUService
	User      *service.UserService
	Meun      *service.MenuService
	Authority *service.AuthorityService
}

type Router struct {
	g      *gin.Engine
	c      *conf.Config
	prefix string
	RouterService
}

func NewRouter(
	c *conf.Config,
	routerSev RouterService,
) *Router {

	r := &Router{
		g: gin.Default(),
		c: c,
		//prefix: "/api",
		RouterService: routerSev,
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
	{
		// 基础路由，提供登录注册
		r.installBaseRouter(publicGroup)
		r.installInitDBRouter(publicGroup)
	}

	// 私有路由
	privateGroup := g.Group(r.prefix)
	// 使用 JWT 中间件和授权中间件
	privateGroup.Use(r.Authmw.AuthFunc(), r.Authzmw.AuthzFunc())

	{
		r.installUserRouter(privateGroup)      // 用户路由
		r.installMenuRouter(privateGroup)      // 菜单路由
		r.installAuthorityRouter(privateGroup) // 角色路由
	}

}

// installBaseRouter 基础路由
func (r *Router) installBaseRouter(group *gin.RouterGroup) {
	baseRouter := group.Group("Base")

	{
		baseRouter.POST("login", r.Base.Login)     // 登录
		baseRouter.POST("captcha", r.Base.Captcha) // 生成验证码
	}
}

// installInitDBRouter 初始化数据库
func (r *Router) installInitDBRouter(group *gin.RouterGroup) {
	initDBRouter := group.Group("init")
	{
		initDBRouter.POST("Initdb", r.Initdb.InitDB)   // 初始化数据库
		initDBRouter.POST("checkdb", r.Initdb.CheckDB) // 检测是否需要初始化数据库
	}
}

// installUserRouter 用户路由
func (r *Router) installUserRouter(group *gin.RouterGroup) {
	userRouter := group.Group("User")
	{
		userRouter.GET("getUserInfo", r.User.GetUserInfo)
		userRouter.POST("changePassword", r.User.ChangePassword)
	}
}

// installMenuRouter 菜单路由
func (r *Router) installMenuRouter(group *gin.RouterGroup) {
	menuRouter := group.Group("menu")
	{
		menuRouter.POST("addBaseMenu", r.Meun.AddBaseMenu)
		menuRouter.POST("addMenuAuthority", r.Meun.AddMenuAuthority) //	增加menu和角色关联关系
		menuRouter.POST("deleteBaseMenu", r.Meun.DeleteBaseMenu)     // 删除菜单
		menuRouter.POST("updateBaseMenu", r.Meun.UpdateBaseMenu)

		menuRouter.POST("getMenu", r.Meun.GetMenu)                   // 获取菜单树
		menuRouter.POST("getMenuList", r.Meun.GetMenuList)           // 分页获取基础menu列表
		menuRouter.POST("getBaseMenuTree", r.Meun.GetBaseMenuTree)   // 获取用户动态路由
		menuRouter.POST("getMenuAuthority", r.Meun.GetMenuAuthority) // 获取指定角色menu
		menuRouter.POST("getBaseMenuById", r.Meun.GetBaseMenuById)   // 根据id获取菜单

	}
}

// installAuthorityRouter 角色路由
func (r *Router) installAuthorityRouter(group *gin.RouterGroup) {
	authorityRouter := group.Group("Authority")

	{
		authorityRouter.POST("getAuthorityList", r.Authority.GetAuthorityList)
	}
}
