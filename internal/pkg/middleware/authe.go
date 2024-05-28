package middleware

import (
	"github.com/gin-gonic/gin"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/response"
	"net/http"
)

// JWTMiddleware JWT 中间件，用于验证刷新 JWT，储存用户基础信息
type JWTMiddleware struct {
	auth auth.Authenticator
}

func NewJWTMiddleware(auth auth.Authenticator) *JWTMiddleware {
	return &JWTMiddleware{
		auth: auth,
	}
}

func (mw *JWTMiddleware) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

// middlewareImpl is the real implementation of the middleware.
func (mw *JWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err := auth.TokenFromHeader(c)
	if err != nil {
		response.ResultWithStatus(c, http.StatusUnauthorized, err)
		c.Abort()
		return
	}

	// parseToken 解析token包含的信息
	userClaims, err := mw.auth.ParseToken(c, token)
	if err != nil {
		response.Result(c, err)
		c.Abort()
		return
	}

	// 将解析出来的用户信息存储到 Gin 的 Context 中
	auth.SetUserClaims(c, userClaims)

	// 刷新token在前端进行，这里不自动刷新
	c.Next()
}
