package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"ktserver/internal/admin/biz"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/model/request"
	"ktserver/internal/pkg/response"
)

// AuthorityService struct
type AuthorityService struct {
	log       *log.Helper
	authority *biz.AuthorityUseCase
}

// NewAuthorityService
func NewAuthorityService(logger log.Logger, authority *biz.AuthorityUseCase) *AuthorityService {
	return &AuthorityService{
		log:       log.NewHelper(logger),
		authority: authority,
	}
}

// GetAuthorityList
// @Tags      Authority
// @Summary   分页获取角色列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=params.AuthorityListResponse,msg=string}  "分页获取角色列表,返回包括列表,总数,页码,每页数量"
// @Router    /authority/getAuthorityList [post]
func (s *AuthorityService) GetAuthorityList(c *gin.Context) {
	var (
		pageInfo request.Pagination
		err      error
		resp     params.AuthorityListResponse
	)

	if err = c.ShouldBindJSON(&pageInfo); err != nil {
		s.log.Error(err)
		response.Result(c, bizerr.InvalidParam.WithMsg(err.Error()))
		return
	}

	resp, err = s.authority.GetAuthorityList(c, pageInfo)

	if err != nil {
		s.log.Error(err)
		response.Result(c, err)
		return
	}

	response.Result(c, resp)
}
