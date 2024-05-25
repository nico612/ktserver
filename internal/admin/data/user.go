package data

import (
	"context"
	"ktserver/internal/admin/model"
)

type UserRepo interface {
	FindUserByName(ctx context.Context, username string) (*model.SysUser, error)
	FindUserByID(ctx context.Context, userID uint) (*model.SysUser, error)
	ChangePassword(userID uint, pwd string) error
}

type userRepo struct {
	data *Data
}

// NewUserRepo .
func NewUserRepo(data *Data) UserRepo {
	return &userRepo{
		data: data,
	}
}

func (r *userRepo) FindUserByName(ctx context.Context, username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.data.DB(ctx).Where("username = ?", username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindUserByID(ctx context.Context, userID uint) (*model.SysUser, error) {
	var user model.SysUser
	if err := r.data.DB(ctx).Where("id = ?", userID).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) ChangePassword(userID uint, pwd string) error {
	return r.data.DB(context.Background()).Model(&model.SysUser{}).Where("id = ?", userID).Update("password", pwd).Error
}
