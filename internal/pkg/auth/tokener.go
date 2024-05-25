package auth

import (
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt/v4"
	"ktserver/internal/pkg/bizerr"
	"ktserver/pkg/utils"
)

type Tokener interface {
	SignToken(claims jwt.Claims) (string, error)
	ParseTokenWithClaims(tokenString string, claims jwt.Claims) (jwt.Claims, error)
}

type JWTOpts struct {
	PrivateKeyFile string
	PublicKeyFile  string
}

func NewDefaultTokener(opts *JWTOpts) (Tokener, error) {
	prikey, err := utils.ReadFile(opts.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	pubkey, err := utils.ReadFile(opts.PublicKeyFile)
	if err != nil {
		return nil, err
	}

	au, err := NewJWTAuth(prikey, pubkey)
	if err != nil {
		return nil, err
	}

	return au, err
}

// esJwtTokn is a jwt auth, 使用ecdsa算法 生成和验证jwt
type esJwtTokn struct {
	issuer string
	prikey *ecdsa.PrivateKey
	pubkey *ecdsa.PublicKey
}

func NewJWTAuth(prikey []byte, pubkey []byte) (Tokener, error) {
	a := &esJwtTokn{}

	if err := a.setPrivateKey(prikey); err != nil {
		return nil, err
	}

	if err := a.setPublicKey(pubkey); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *esJwtTokn) setPrivateKey(prikey []byte) error {
	pk, err := jwt.ParseECPrivateKeyFromPEM(prikey)
	if err != nil {
		return err
	}
	a.prikey = pk
	return nil
}

func (a *esJwtTokn) setPublicKey(pubkey []byte) error {
	pk, err := jwt.ParseECPublicKeyFromPEM(pubkey)
	if err != nil {
		return err
	}
	a.pubkey = pk
	return nil
}

// SignToken 生成jwt token
func (a *esJwtTokn) SignToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ts, err := token.SignedString(a.prikey)
	if err != nil {
		return "", bizerr.ErrSignTokenFailed
	}
	return ts, nil
}

// ParseToken 解析jwt token
func (a *esJwtTokn) ParseTokenWithClaims(tokenString string, claims jwt.Claims) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return a.pubkey, nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
