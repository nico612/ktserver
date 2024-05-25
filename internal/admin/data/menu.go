package data

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/model"
)

type MenuRepo interface {
	// UserAuthorityDefaultRouter 获取用户默认路由
	UserAuthorityDefaultRouter(ctx context.Context, defaultAuthorityRouter string, authorityId uint) (string, error)
}

const NOT_FOUND_ROUTER = "404"

type menuRepo struct {
	data *Data
}

func NewMenuRepo(data *Data) MenuRepo {
	return &menuRepo{
		data: data,
	}
}

// UserAuthorityDefaultRouter 检查用户默认路由
func (r *menuRepo) UserAuthorityDefaultRouter(ctx context.Context, defaultAuthorityRouter string, authorityId uint) (string, error) {

	// 根据用户角色获取菜单ids
	var menuIds []string
	err := r.data.DB(ctx).Model(&model.SysAuthorityMenu{}).Where("sys_authority_authority_id = ?", authorityId).Pluck("sys_base_menu_id", &menuIds).Error
	if err != nil {
		return defaultAuthorityRouter, err
	}

	// 根据用户默认路由获取菜单
	var am model.SysBaseMenu
	err = r.data.DB(ctx).Where("name = ? and id in (?)", defaultAuthorityRouter, menuIds).First(&am).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return NOT_FOUND_ROUTER, nil
	}

	return defaultAuthorityRouter, err
}
