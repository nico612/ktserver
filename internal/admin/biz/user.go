package biz

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/data"
	"ktserver/internal/admin/model"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/lock"
	"ktserver/pkg/utils"
)

type UserUseCase struct {
	userRepo data.UserRepo
	locker   *lock.RedisLocker
}

func NewUserUseCase(userRepo data.UserRepo, locker *lock.RedisLocker) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		locker:   locker,
	}
}

func (uc *UserUseCase) FindUserByID(ctx context.Context, userID uint) (*model.SysUser, error) {
	return uc.userRepo.FindUserByID(ctx, userID)
}

func (uc *UserUseCase) ChangePassword(c context.Context, userID uint, req params.ChangePasswordReq) error {
	user, err := uc.FindUserByID(c, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.UserNotFound
		}
		return err
	}
	if !utils.BcryptCheck(req.Password, user.Password) {
		return bizerr.InvalidPassword

	}
	newPwd := utils.BcryptHash(req.NewPassword)
	if err = uc.userRepo.ChangePassword(userID, newPwd); err != nil {
		return err
	}
	return nil
}
