package data

import (
	"context"
	"ktserver/internal/admin/model"
	"ktserver/internal/pkg/model/request"
)

// AuthorityRepo
type AuthorityRepo interface {
	GetAuthorityList(ctx context.Context, pageInfo request.Pagination) (total int64, list []model.SysAuthority, err error)
}

type authorityRepo struct {
	data *Data
}

// NewAuthorityRepo
func NewAuthorityRepo(data *Data) AuthorityRepo {
	return &authorityRepo{data: data}
}

func (r *authorityRepo) GetAuthorityList(ctx context.Context, pageInfo request.Pagination) (total int64, list []model.SysAuthority, err error) {

	db := r.data.DB(ctx).Model(&model.SysAuthority{}).Where("parent_id = ?", 0)

	if err = db.Count(&total).Error; err != nil || total == 0 {
		return 0, nil, err
	}

	if err = db.Limit(pageInfo.Limit).Offset(pageInfo.Offset).Preload("DataAuthorityId").Find(&list).Error; err != nil {
		return 0, nil, err
	}

	return
}
