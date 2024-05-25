package auth

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type UserClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
	ExpireTime  time.Duration
}

// Custom claims structure
type JWTClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

func NewUserClaims(userClaims UserClaims) *JWTClaims {
	return &JWTClaims{
		UserClaims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ktserver",                                                // 签名的发行者
			Subject:   strconv.Itoa(int(userClaims.ID)),                          // 主题, 可以用来鉴别一个用户
			Audience:  jwt.ClaimStrings{"ktserver"},                              // 受众, 一般可以为特定的App，服务或模块。服务端的安全策略在签发时和验证时，aud必须时一致的；
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(userClaims.ExpireTime)), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),      // 生效时间，这里提前一秒， 避免没有生效的问题
			IssuedAt:  jwt.NewNumericDate(time.Now()),                            // 签发时间
			ID:        "",                                                        // 令牌唯一标识符，通常用于一次性消费的Token
		},
	}
}
