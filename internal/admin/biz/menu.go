package biz

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/data"
	"ktserver/internal/admin/model"
	"ktserver/internal/pkg/bizerr"
)

type MenuUseCase struct {
	menuRepo data.MenuRepo
}

func NewMenuUseCase(menuRepo data.MenuRepo) *MenuUseCase {
	return &MenuUseCase{
		menuRepo: menuRepo,
	}
}

func (uc *MenuUseCase) GetMenuTree(ctx context.Context, authorityId uint) (menus []model.SysMenu, err error) {

	// 1. 根据角色id获取所有的菜单树
	menuTreeMap, err := uc.menuRepo.GetMenuTreeMap(ctx, authorityId)
	if err != nil {
		return nil, err
	}

	// 2. 获取根菜单列表
	menus = menuTreeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = uc.getChildrenList(ctx, &menus[i], menuTreeMap)
	}

	return menus, err

}

func (uc *MenuUseCase) getChildrenList(ctx context.Context, menu *model.SysMenu, menuTreeMap map[string][]model.SysMenu) (err error) {
	menu.Children = menuTreeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = uc.getChildrenList(ctx, &menu.Children[i], menuTreeMap)
	}
	return err
}

func (uc *MenuUseCase) AddBaseMenu(c context.Context, menu model.SysBaseMenu) error {
	if _, err := uc.menuRepo.FindBaseMenuByName(c, menu.Name); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uc.menuRepo.AddBaseMenu(c, menu)
		}
		return err
	}
	return bizerr.MenuExist
}
