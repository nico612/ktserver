package data

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/model"
	"strconv"
)

type MenuRepo interface {
	// UserAuthorityDefaultRouter 获取用户默认路由
	UserAuthorityDefaultRouter(ctx context.Context, defaultAuthorityRouter string, authorityId uint) (string, error)
	// GetMenuTreeMap 根据角色id获取菜单树
	GetMenuTreeMap(ctx context.Context, authorityId uint) (treeMap map[string][]model.SysMenu, err error)
	AddBaseMenu(ctx context.Context, menu model.SysBaseMenu) error
	FindBaseMenuByName(ctx context.Context, name string) (model.SysBaseMenu, error)
}

const NOT_FOUND_ROUTER = "404"

type menuRepo struct {
	data *datastore
}

func newMenuRepo(data *datastore) MenuRepo {
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

func (r *menuRepo) GetMenuTreeMap(ctx context.Context, authorityId uint) (treeMap map[string][]model.SysMenu, err error) {
	var allMenus []model.SysMenu
	var baseMenus []model.SysBaseMenu
	var btns []model.SysAuthorityBtn

	treeMap = make(map[string][]model.SysMenu)

	// 1. 根据角色id获取所有的授权菜单
	var sysAuthorityMenus []model.SysAuthorityMenu
	if err = r.data.DB(ctx).Where("sys_authority_authority_id = ?", authorityId).Find(&sysAuthorityMenus).Error; err != nil {
		return nil, err
	}

	// 2. 根据授权菜单获取对应的菜单信息
	var menuIds []string
	for i := range sysAuthorityMenus {
		menuIds = append(menuIds, sysAuthorityMenus[i].MenuId)
	}

	if err = r.data.DB(ctx).Where("id in (?)", menuIds).Order("sort").Preload("Parameters").Find(&baseMenus).Error; err != nil {
		return nil, err
	}

	// 3. 将菜单信息转换为SysMenu模型
	for i := range baseMenus {
		allMenus = append(allMenus, model.SysMenu{
			SysBaseMenu: baseMenus[i],
			AuthorityId: authorityId,
			MenuId:      strconv.Itoa(int(baseMenus[i].ID)),
			Parameters:  baseMenus[i].Parameters,
		})
	}

	// 4. 获取所有的按钮
	if err = r.data.DB(ctx).Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error; err != nil {
		return nil, err
	}

	// 5. 将按钮信息转换为SysMenu模型
	// sysMenuID => btnName => authorityId
	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if _, ok := btnMap[v.SysMenuID]; !ok {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}

	// 6. 构建菜单树
	for _, v := range allMenus {
		v.Btns = btnMap[v.SysBaseMenu.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}

	return treeMap, err
}

func (r *menuRepo) AddBaseMenu(ctx context.Context, menu model.SysBaseMenu) error {

	return r.data.DB(ctx).Create(&menu).Error
}

func (r *menuRepo) FindBaseMenuByName(ctx context.Context, name string) (model.SysBaseMenu, error) {
	var menu model.SysBaseMenu
	err := r.data.DB(ctx).Where("name = ?", name).First(&menu).Error
	return menu, err
}
