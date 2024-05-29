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
	NewStore,
	wire.Bind(new(IStore), new(*datastore)),
)

// transactionKey is an unique key used in context to store
// transaction instances to be shared between multiple operations.
type transactionKey struct{}

var (
	once sync.Once
	S    *datastore
)

type IStore interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) (err error)
	Users() UserRepo
	Menus() MenuRepo
	Authorities() AuthorityRepo
}

var _ IStore = (*datastore)(nil)

// datastore .
type datastore struct {
	db *gorm.DB
}

// NewStore .
func NewStore(db *gorm.DB, logger log.Logger) (*datastore, func(), error) {
	once.Do(func() {
		S = &datastore{db: db}
	})
	return S, cleanup, nil
}

func (ds *datastore) Users() UserRepo {
	return newUserRepo(ds)
}

func (ds *datastore) Menus() MenuRepo {
	return newMenuRepo(ds)
}

func (ds *datastore) Authorities() AuthorityRepo {
	return newAuthorityRepo(ds)
}

// DB .从context中获取事务，如果不存在则返回db
func (ds *datastore) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return ds.db
}

// Tx 事务
func (ds *datastore) Tx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务存入context
		ctx = context.WithValue(ctx, transactionKey{}, tx)
		return fn(ctx)
	})
}

// AutoMigrate 数据库迁移
func AutoMigrate() error {
	if S == nil {
		return nil
	}
	return S.db.AutoMigrate(
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
func cleanup() {
	sqlDB, _ := S.db.DB()
	sqlDB.Close()
}
