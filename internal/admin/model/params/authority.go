package params

import (
	"ktserver/internal/admin/model"
	"ktserver/internal/pkg/model/request"
)

// 角色列表响应
type AuthorityListResponse struct {
	List []model.SysAuthority `json:"list"`
	request.Pagination
}
