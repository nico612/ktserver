package source

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ktserver/internal/admin/model"
)

type initDictDetail struct{}

func (i *initDictDetail) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&model.SysDictionaryDetail{})
}

func (i *initDictDetail) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&model.SysDictionaryDetail{})
}

func (i initDictDetail) InitializerName() string {
	return model.SysDictionaryDetail{}.TableName()
}

func (i *initDictDetail) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, ErrMissingDBContext
	}
	dicts, ok := ctx.Value(initDict{}.InitializerName()).([]model.SysDictionary)
	if !ok {
		return ctx, errors.Wrap(ErrMissingDependentContext,
			fmt.Sprintf("未找到 %s 表初始化数据", model.SysDictionary{}.TableName()))
	}
	True := true
	dicts[0].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "男", Value: "1", Status: &True, Sort: 1},
		{Label: "女", Value: "2", Status: &True, Sort: 2},
	}

	dicts[1].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "smallint", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "mediumint", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "int", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "bigint", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "int2", Value: "5", Status: &True, Extend: "pgsql", Sort: 5},
		{Label: "int4", Value: "6", Status: &True, Extend: "pgsql", Sort: 6},
		{Label: "int6", Value: "7", Status: &True, Extend: "pgsql", Sort: 7},
		{Label: "int8", Value: "8", Status: &True, Extend: "pgsql", Sort: 8},
	}

	dicts[2].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "date", Status: &True},
		{Label: "time", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "year", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "datetime", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "timestamp", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "timestamptz", Value: "6", Status: &True, Extend: "pgsql", Sort: 5},
	}
	dicts[3].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "float", Status: &True},
		{Label: "double", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "decimal", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "numeric", Value: "3", Status: &True, Extend: "pgsql", Sort: 3},
		{Label: "smallserial", Value: "4", Status: &True, Extend: "pgsql", Sort: 4},
	}

	dicts[4].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "char", Status: &True},
		{Label: "varchar", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "tinyblob", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "tinytext", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "text", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "blob", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "mediumblob", Value: "6", Status: &True, Extend: "mysql", Sort: 6},
		{Label: "mediumtext", Value: "7", Status: &True, Extend: "mysql", Sort: 7},
		{Label: "longblob", Value: "8", Status: &True, Extend: "mysql", Sort: 8},
		{Label: "longtext", Value: "9", Status: &True, Extend: "mysql", Sort: 9},
	}

	dicts[5].SysDictionaryDetails = []model.SysDictionaryDetail{
		{Label: "tinyint", Value: "1", Extend: "mysql", Status: &True},
		{Label: "bool", Value: "2", Extend: "pgsql", Status: &True},
	}
	for _, dict := range dicts {
		if err := db.Model(&dict).Association("SysDictionaryDetails").
			Replace(dict.SysDictionaryDetails); err != nil {
			return ctx, errors.Wrap(err, model.SysDictionaryDetail{}.TableName()+"表数据初始化失败!")
		}
	}
	return ctx, nil
}

func (i *initDictDetail) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var dict model.SysDictionary
	if err := db.Preload("SysDictionaryDetails").
		First(&dict, &model.SysDictionary{Name: "数据库bool类型"}).Error; err != nil {
		return false
	}
	return len(dict.SysDictionaryDetails) > 0 && dict.SysDictionaryDetails[0].Label == "tinyint"
}
