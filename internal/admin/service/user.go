package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"ktserver/internal/admin/biz"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/auth"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/response"
)

type UserService struct {
	c           *conf.Config
	log         *log.Helper
	userUseCase *biz.UserUseCase
}

func NewUserService(c *conf.Config, logger log.Logger, userUserCase *biz.UserUseCase) *UserService {
	return &UserService{
		c:           c,
		log:         log.NewHelper(logger),
		userUseCase: userUserCase,
	}
}

// GetUserInfo
// @Tags      SysUser
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取用户信息"
// @Router    /user/getUserInfo [get]
func (s *UserService) GetUserInfo(c *gin.Context) {
	userID := auth.UserIDFromGinCtx(c)
	user, err := s.userUseCase.FindUserByID(c, userID)
	if err != nil {
		response.Result(c, err)
		return
	}
	response.Result(c, params.UserInfoResp{UserInfo: lo.FromPtr(user)})
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      parma.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  response.Response{msg=string}  "用户修改密码"
// @Router    /user/changePassword [post]
func (s *UserService) ChangePassword(c *gin.Context) {
	var req params.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Result(c, bizerr.InvalidParam.WithMsg(err.Error()))
		return
	}
	userID := auth.UserIDFromGinCtx(c)
	err := s.userUseCase.ChangePassword(c, userID, req)

	response.Result(c, err)
}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
