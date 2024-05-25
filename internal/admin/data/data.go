package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
	"ktserver/internal/admin/model"
	"sync"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewMenuRepo,
	NewAuthorityRepo,
)

// transactionKey is an unique key used in context to store
// transaction instances to be shared between multiple operations.
type transactionKey struct{}

var (
	once sync.Once
	DS   *Data
)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	once.Do(func() {
		DS = &Data{db: db}
	})

	cleanup := func() {
		Close()
		log.NewHelper(logger).Info("closed the data resources")
	}
	return DS, cleanup, nil
}

// DB .从context中获取事务，如果不存在则返回db
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

// Tx 事务
func (d *Data) Tx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务存入context
		ctx = context.WithValue(ctx, transactionKey{}, tx)
		return fn(ctx)
	})
}

// AutoMigrate 数据库迁移
func AutoMigrate() error {
	if DS == nil {
		return nil
	}
	return DS.db.AutoMigrate(
		&model.SysApi{},
		&model.SysUser{},
		&model.SysBaseMenu{},
		&model.SysAuthority{},
		&model.SysDictionary{},
		&model.SysBaseMenuParameter{},
		&model.SysBaseMenuBtn{},
		&model.SysAuthorityBtn{},
	)
}

// Close gorm.db
func Close() {
	sqlDB, _ := DS.db.DB()
	sqlDB.Close()
}
