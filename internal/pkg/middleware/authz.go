package middleware

import (
	"github.com/gin-gonic/gin"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/authz"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/response"
	"strconv"
	"strings"
)

// AuthzMiddleware 权限中间件，用于验证用户权限
type AuthzMiddleware struct {
	authz *authz.CasbinAuthorizer
}

func NewAuthzMiddleware(authz *authz.CasbinAuthorizer) *AuthzMiddleware {
	return &AuthzMiddleware{
		authz: authz,
	}
}

func (mw *AuthzMiddleware) AuthzFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *AuthzMiddleware) middlewareImpl(c *gin.Context) {
	//获取请求的PATH
	path := c.Request.URL.Path
	// 去掉前缀
	//obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
	obj := strings.TrimSpace(path)

	// 获取请求方法
	act := c.Request.Method

	// 获取用户的角色
	claims, err := auth.UserClaimsFromGinCtx(c)
	if err != nil {
		response.Result(c, err)
		c.Abort()
		return
	}

	// 这里要将int64转为string 否则查询授权会失败
	sub := strconv.Itoa(int(claims.AuthorityId))

	// 判断策略中是否存在
	success, _ := mw.authz.Enforcer.Enforce(sub, obj, act)
	if !success {
		response.Result(c, bizerr.ErrNoPermission)
		c.Abort()
		return
	}

	c.Next()
}
