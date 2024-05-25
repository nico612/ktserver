package handler

import (
	"context"
	"fmt"
	"ktserver/internal/pkg/initdb"
	"ktserver/pkg/db"
)

type MysqlInitHandler struct{}

var _ initdb.TypedDBInitHandler = (*MysqlInitHandler)(nil)

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

// EnsureDB 创建数据库并初始化 mysql
func (h *MysqlInitHandler) EnsureDB(ctx context.Context, conf *initdb.InitDBModel) (context.Context, error) {
	if conf.DBName == "" {
		return ctx, nil
	}

	mysqlOpts := conf.ToMysqlConfig()

	// 创建数据库
	dsn := conf.MysqlEmptyDsn()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.DBName)
	if err := createDatabase(dsn, "mysql", createSql); err != nil {
		return nil, err
	}
	// 初始化数DB
	gormDB, err := db.NewMySQL(mysqlOpts)
	if err != nil {
		return ctx, nil
	}

	return context.WithValue(ctx, "db", gormDB), nil
}

func (h *MysqlInitHandler) InitTables(ctx context.Context, inits []initdb.SubInitializer) error {
	return createTables(ctx, inits)
}

func (h *MysqlInitHandler) InitData(ctx context.Context, inits []initdb.SubInitializer) error {
	next, cancle := context.WithCancel(ctx)
	defer cancle()

	for _, init := range inits {
		if init.DataInserted(next) {
			continue
		}

		if n, err := init.InitializeData(next); err != nil {
			return err
		} else {
			next = n
		}
	}
	return nil
}

//// WriteConfig mysql回写配置
//func (h *MysqlInitHandler) WriteConfig(ctx context.Context) error {
//	c, ok := ctx.Value("config").(conf.Config)
//	if !ok {
//		return errors.New("mysql config invalid")
//	}
//
//	fmt.Printf("mysql config: %+v\n", c)
//
//	return nil
//}
