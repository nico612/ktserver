package biz

import (
	"context"
	"ktserver/internal/admin/data"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/model/request"
)

// AuthorityUseCase
type AuthorityUseCase struct {
	store data.IStore
}

// NewAuthorityUseCase
func NewAuthorityUseCase(store data.IStore) *AuthorityUseCase {
	return &AuthorityUseCase{
		store: store,
	}
}

// GetAuthorityList 获取角色列表
func (c *AuthorityUseCase) GetAuthorityList(ctx context.Context, pageInfo request.Pagination) (params.AuthorityListResponse, error) {
	pageInfo.Check()
	total, list, err := c.store.Authorities().GetAuthorityList(ctx, pageInfo)
	if err != nil {
		return params.AuthorityListResponse{}, err
	}

	pageInfo.Total = total
	return params.AuthorityListResponse{
		List:       list,
		Pagination: pageInfo,
	}, nil
}
