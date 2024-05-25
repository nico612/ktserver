package initdb

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"sort"
)

/* ---- * biz * ---- */

type InitDBService struct {
	initHandler TypedDBInitHandler
}

func NewInitDBService(initHandler TypedDBInitHandler) *InitDBService {
	return &InitDBService{
		initHandler: initHandler,
	}
}

// InitDBModel 创建数据库并初始化 总入口
func (s *InitDBService) InitDB(conf *InitDBModel) (err error) {
	ctx := context.TODO()
	if len(initializers) == 0 {
		return errors.New("无可用初始化过程，请检查初始化是否已执行完成")
	}
	sort.Sort(&initializers) // 保证有依赖的 initializer 排在后面执行

	// Note: 若 initializer 只有单一依赖，可以写为 B=A+1, C=A+1; 由于 BC 之间没有依赖关系，所以谁先谁后并不影响初始化
	// 若存在多个依赖，可以写为 C=A+B, D=A+B+C, E=A+1;
	// C必然>A|B，因此在AB之后执行，D必然>A|B|C，因此在ABC后执行，而E只依赖A，顺序与CD无关，因此E与CD哪个先执行并不影响
	//var initHandler TypedDBInitHandler

	switch conf.DBType {
	case Mysql:
		//initHandler = handler.NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", Mysql)
	//case Pgsql:
	//	initHandler = NewPgsqlInitHandler()
	//	ctx = context.WithValue(ctx, "dbtype", Pgsql)
	//case Sqlite:
	//	initHandler = NewSqliteInitHandler()
	//	ctx = context.WithValue(ctx, "dbtype", Sqlite)
	//case Mssql:
	//	initHandler = NewMssqlInitHandler()
	//	ctx = context.WithValue(ctx, "dbtype", Mssql)
	default:
		return errors.New("不支持的数据库类型")
	}

	// 创建数据库
	ctx, err = s.initHandler.EnsureDB(ctx, conf)
	if err != nil {
		return err
	}

	subInitializer := lo.Map(initializers, func(i *orderedInitializer, index int) SubInitializer {
		return i.SubInitializer
	})

	// 初始化表
	if err = s.initHandler.InitTables(ctx, subInitializer); err != nil {
		return err
	}

	// 初始化数据
	if err = s.initHandler.InitData(ctx, subInitializer); err != nil {
		return err
	}

	//// 回写配置
	//if err = initHandler.WriteConfig(ctx); err != nil {
	//	return err
	//}

	// 清空初始化器

	initializers = initSlice{}
	cache = map[string]*orderedInitializer{}

	return nil

}
