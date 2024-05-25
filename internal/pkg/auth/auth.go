package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/genericoptions"
	"time"
)

type Authenticator interface {
	NewToken(ctx context.Context, userClaims UserClaims) (token string, err error)
	ParseToken(ctx context.Context, token string) (claims UserClaims, err error)
	IsBlocked(ctx context.Context, token string) (blocked bool, err error)
	BlockToken(ctx context.Context, token string, expireTime time.Time) error
	CreateTokenByOldToken(ctx context.Context, oldToken string, userClaims UserClaims) (token string, err error)
}

type ctxClaims struct{}

var (
	TokenIndexKey = "x-token"
	CtxClaimsKey  = "ctx-claims-key"
)

type JWTAuthenticator struct {
	tokener  Tokener               // 生成和解析token的方法
	rds      redis.UniversalClient // redis client 用于存储token的block状态
	group    *singleflight.Group   // 旧token 换新token 使用归并回源避免并发问题
	authOpts *genericoptions.AuthOptions
}

func blockKey(token string) string {
	return fmt.Sprintf("auth:block:%x", md5.Sum([]byte(token)))
}

func cacheKey(userID uint) string {
	return fmt.Sprintf("auth:cache:%d", userID)
}

func NewAuthenticator(tokener Tokener, authOpts *genericoptions.AuthOptions, rds redis.UniversalClient) Authenticator {
	return &JWTAuthenticator{
		tokener:  tokener,
		authOpts: authOpts,
		rds:      rds,
		group:    &singleflight.Group{},
	}
}

func (a *JWTAuthenticator) NewToken(ctx context.Context, userClaims UserClaims) (token string, err error) {
	claims := NewUserClaims(userClaims)
	token, err = a.tokener.SignToken(claims)
	if err != nil {
		return
	}

	// 存储token
	if err = a.AddToken(ctx, token, claims.UserClaims, userClaims.ExpireTime); err != nil {
		return
	}

	return token, nil
}

// ParseToken 解析jwt token
func (a *JWTAuthenticator) ParseToken(ctx context.Context, token string) (userClaims UserClaims, err error) {

	// 解析token
	userClaims, err = a.parseToken(token)
	if err != nil {
		return userClaims, err
	}

	// 验证token是否被锁定
	blocked, err := a.IsBlocked(ctx, token)
	if err != nil {
		return userClaims, err
	}
	if blocked {
		return userClaims, bizerr.ErrTokenInvalid
	}

	// 从缓存中读取token, 判断当前token是否有效
	tokenStr, err := a.GetToken(ctx, userClaims)
	if err != nil {
		return
	}
	if len(tokenStr) == 0 {
		return userClaims, bizerr.ErrTokenExpired
	}

	if token != tokenStr { // token已失效，在新的地方登录，产生新的token
		return userClaims, bizerr.ErrUserLoginOtherDevice
	}

	return userClaims, nil
}

// parseToken 解析token
func (a *JWTAuthenticator) parseToken(tokenString string) (userClaims UserClaims, err error) {

	claims, err := a.tokener.ParseTokenWithClaims(tokenString, &JWTClaims{})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok {
			return userClaims, bizerr.ErrTokenParseFail
		}
		if ve.Errors&jwt.ValidationErrorMalformed != 0 { // 无效的token
			return userClaims, bizerr.ErrTokenInvalid
		}

		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// token已过期 或者 token未生效， 这里返回token过期错误
			return userClaims, bizerr.ErrTokenExpired
		}
		// 其他错误
		return userClaims, bizerr.ErrTokenParseFail
	}

	jwtClaims, ok := claims.(*JWTClaims)
	if !ok {
		return userClaims, bizerr.ErrTokenParseFail
	}

	//// 验证issuer, 其实这里不需要验证，因为使用公私钥签名，只有私钥持有者才能签名，公钥持有者才能解密
	//if jwtClaims.Issuer != a.authOpts.Issuer {
	//	return userClaims, bizerr.ErrTokenInvalid
	//}
	//
	//// 验证uid
	//uid, err := strconv.Atoi(jwtClaims.Subject)
	//if err != nil || uid == 0 {
	//	return userClaims, bizerr.ErrTokenInvalid
	//}

	return jwtClaims.UserClaims, nil
}

func (a *JWTAuthenticator) CreateTokenByOldToken(ctx context.Context, oldToken string, userClaims UserClaims) (token string, err error) {
	// 生成新token
	v, err, _ := a.group.Do("JWT:"+oldToken, func() (interface{}, error) {
		// 生成新token
		return a.NewToken(ctx, userClaims)
	})
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

// IsBlocked 判断token是否被锁定, 用于限制同平台登录，该项目未使用到
func (a *JWTAuthenticator) IsBlocked(ctx context.Context, token string) (blocked bool, err error) {
	key := blockKey(token)

	exists, err := a.rds.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// BlockToken 锁定token
// expireTime：原来token的过期时间即可
func (a *JWTAuthenticator) BlockToken(ctx context.Context, token string, expireTime time.Time) error {
	dur := expireTime.Sub(time.Now())
	if dur > 0 {
		key := blockKey(token)
		if err := a.rds.SetNX(ctx, key, "", dur).Err(); err != nil {
			return err
		}
	}
	return nil
}

// AddToken 存储token to redis
func (a *JWTAuthenticator) AddToken(ctx context.Context, token string, userClaims UserClaims, expireTime time.Duration) error {
	key := cacheKey(userClaims.ID)
	return a.rds.Set(ctx, key, token, expireTime).Err()
}

func (a *JWTAuthenticator) GetToken(ctx context.Context, userClaims UserClaims) (string, error) {
	key := cacheKey(userClaims.ID)
	return a.rds.Get(ctx, key).Result()
}

func UserClaimsFromGinCtx(c *gin.Context) (UserClaims, error) {
	v, exists := c.Get(CtxClaimsKey)
	if !exists {
		return UserClaims{}, bizerr.ErrTokenParseFail
	}
	return v.(UserClaims), nil
}

func UserIDFromGinCtx(c *gin.Context) uint {
	claims, err := UserClaimsFromGinCtx(c)
	if err != nil {
		return 0
	}
	return claims.ID
}

func SetUserClaims(c *gin.Context, userClaims UserClaims) {
	c.Set(CtxClaimsKey, userClaims)
}

func TokenFromHeader(c *gin.Context) (string, error) {

	token := c.Request.Header.Get(TokenIndexKey)

	if len(token) == 0 {
		return "", bizerr.UnAuthority
	}

	return token, nil
}

//func WithClaims(ctx context.Context, usrClaims UserClaims) context.Context {
//	return context.WithValue(ctx, CtxClaimsKey, usrClaims)
//}
//
//func ClaimsFromCtx(ctx context.Context) (UserClaims, error) {
//	claims, ok := ctx.Value(CtxClaimsKey).(JWTClaims)
//	if !ok {
//		return UserClaims{}, bizerr.UnAuthority
//	}
//	return claims.UserClaims, nil
//}
//
//func UserIDFromGinCtx(ctx context.Context) (userID uint) {
//	userClaims, _ := ClaimsFromCtx(ctx)
//	return userClaims.ID
//}
//
