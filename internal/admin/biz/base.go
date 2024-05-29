package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/data"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/bizerr"
	"ktserver/pkg/utils"
	"time"
)

type BaseUseCase struct {
	c     *conf.Config
	log   *log.Helper
	store data.IStore
	//menuRepo      data.MenuRepo
	authenticator auth.Authenticator
}

func NewBaseUseCase(
	c *conf.Config,
	logger log.Logger,
	authenticator auth.Authenticator,
	store data.IStore,
	// userRepo data.UserRepo,
	// menuRepo data.MenuRepo,
) (*BaseUseCase, error) {

	return &BaseUseCase{
		c:             c,
		log:           log.NewHelper(logger),
		authenticator: authenticator,
		store:         store,
	}, nil
}

func (uc *BaseUseCase) Login(c context.Context, req params.Login) (resp params.LoginResponse, err error) {
	user, err := uc.store.Users().FindUserByName(c, req.Username)
	if err != nil {
		return
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		err = bizerr.LoginPasswordError
		return
	}
	// userAuthority
	userDefaultRouter, err := uc.store.Menus().UserAuthorityDefaultRouter(c, user.Authority.DefaultRouter, user.AuthorityId)
	if err != nil {
		return
	}

	user.Authority.DefaultRouter = userDefaultRouter

	// jwt token
	userClaims := auth.UserClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Username:    user.Username,
		NickName:    user.NickName,
		AuthorityId: user.AuthorityId,
		ExpireTime:  uc.c.AuthOptions.ExpireTime,
	}

	resp.ExpiresAt = time.Now().Add(userClaims.ExpireTime).Unix()
	resp.User = *user
	token, err := uc.authenticator.NewToken(c, userClaims)
	if err != nil {
		return
	}
	resp.Token = token

	return
}
