package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"ktserver/internal/admin/biz"
	"ktserver/internal/admin/model"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/response"
)

type MenuService struct {
	log *log.Helper

	menu *biz.MenuUseCase
}

func NewMenuService(logger log.Logger, menu *biz.MenuUseCase) *MenuService {
	return &MenuService{
		log:  log.NewHelper(logger),
		menu: menu,
	}
}

// AddBaseMenu
// @Tags      Menu
// @Summary   新增菜单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      params.AddBaseMenuReq             true  "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success   200   {object}  response.Response{msg=string}  "新增菜单"
// @Router    /menu/addBaseMenu [post]
func (s *MenuService) AddBaseMenu(c *gin.Context) {
	var menu params.AddBaseMenuReq
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.Result(c, bizerr.InvalidParam)
		return
	}

	err = s.menu.AddBaseMenu(c, menu.SysBaseMenu)
	if err != nil {
		s.log.Error(err)
		response.Result(c, err)
	}

	response.Result(c, nil)
}

func (s *MenuService) AddMenuAuthority(c *gin.Context) {

}

func (s *MenuService) DeleteBaseMenu(c *gin.Context) {

}

func (s *MenuService) UpdateBaseMenu(c *gin.Context) {

}

// GetMenu
// @Tags      AuthorityMenu
// @Summary   获取用户动态路由
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.Empty                                                  true  "空"
// @Success   200   {object}  response.Response{data=params.SysMenusResponse,msg=string}  "获取用户动态路由,返回包括系统菜单详情列表"
// @Router    /menu/getMenu [post]
func (s *MenuService) GetMenu(c *gin.Context) {

	userClaims, err := auth.UserClaimsFromGinCtx(c)
	if err != nil {
		s.log.Error(err)
		response.Result(c, err)
	}

	menus, err := s.menu.GetMenuTree(c, userClaims.AuthorityId)
	if err != nil {
		s.log.Error(err)
		response.Result(c, err)
	}

	if menus == nil {
		menus = []model.SysMenu{}
	}

	resp := params.SysMenusResponse{Menus: menus}

	response.Result(c, resp)
}

func (s *MenuService) GetMenuList(c *gin.Context) {

}

func (s *MenuService) GetBaseMenuTree(c *gin.Context) {

}

func (s *MenuService) GetMenuAuthority(c *gin.Context) {

}

func (s *MenuService) GetBaseMenuById(c *gin.Context) {

}
