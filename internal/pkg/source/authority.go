package source

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/model"
	"ktserver/internal/pkg/initdb"
	"ktserver/pkg/utils"
)

// 角色数据

var _ initdb.SubInitializer = (*initAuthority)(nil)

type initAuthority struct{}

func (i initAuthority) InitializerName() string {
	return model.SysAuthority{}.TableName()
}

func (i *initAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&model.SysAuthority{})
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, ErrMissingDBContext
	}

	// 初始化角色数据
	entities := []model.SysAuthority{
		{AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "测试角色", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "普通用户子角色", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", model.SysAuthority{}.TableName())
	}

	// data authority
	if err := db.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*model.SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		},
	); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			db.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	if err := db.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*model.SysAuthority{
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			db.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)

	return next, nil
}

func (i *initAuthority) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	return db.Migrator().HasTable(&model.SysAuthority{})
}

func (i *initAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	if errors.Is(db.Where("authority_id = ?", "8881").
		First(&model.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true

}
