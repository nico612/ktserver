package initdb

import (
	"context"
	"errors"
	"fmt"
)

const (
	Mysql           = "mysql"
	Pgsql           = "pgsql"
	Sqlite          = "sqlite"
	Mssql           = "mssql"
	InitSuccess     = "\n[%v] --> 初始数据成功!\n"
	InitDataExist   = "\n[%v] --> %v 的初始数据已存在!\n"
	InitDataFailed  = "\n[%v] --> %v 初始数据失败! \nerr: %+v\n"
	InitDataSuccess = "\n[%v] --> %v 初始数据成功!\n"
)

const (
	InitOrderInternal = 1000
	InitOrderExternal = 100000
)

var (
	ErrMissingDependentContext = errors.New("missing dependent value in context")
	ErrDBTypeMismatch          = errors.New("db type mismatch")
)

// SubInitializer 表初始化器 提供 source/*/init() 使用的接口，每个 initializer 完成一个初始化过程
type SubInitializer interface {
	InitializerName() string // 表名，不一定代表单独一个表，所以改成了更宽泛的语义
	MigrateTable(ctx context.Context) (next context.Context, err error)
	InitializeData(ctx context.Context) (next context.Context, err error)
	TableCreated(ctx context.Context) bool
	DataInserted(ctx context.Context) bool
}

// TypedDBInitHandler 执行传入的 initializer 数据库初始化处理器
type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context, conf *InitDBModel) (context.Context, error) // 建库，失败属于 fatal error，因此让它 panic
	//WriteConfig(ctx context.Context) error                               // 回写配置
	InitTables(ctx context.Context, inits []SubInitializer) error // 建表 handler
	InitData(ctx context.Context, inits []SubInitializer) error   // 建数据 handler
}

// orderedInitializer 组合一个顺序字段，以供排序
type orderedInitializer struct {
	order int
	SubInitializer
}

// initSlice 供 initializer 排序依赖时使用
type initSlice []*orderedInitializer

// 实现排序方法

func (s initSlice) Len() int {
	return len(s)
}

func (s initSlice) Less(i, j int) bool {
	return s[i].order < s[j].order
}

func (s initSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var (
	initializers initSlice // 储存所有的初始化器，需通过 RegisterInit 注册
	cache        map[string]*orderedInitializer
)

// RegisterInit 注册要执行的初始化过程，会在 InitDBModel() 时调用
// order: 顺序
// i: 初始化器
func RegisterInit(order int, i SubInitializer) {
	if initializers == nil {
		initializers = initSlice{}
	}
	if cache == nil {
		cache = map[string]*orderedInitializer{}
	}

	// 表名
	name := i.InitializerName()
	if _, existed := cache[name]; existed {
		panic(fmt.Sprintf("Name conflict on %s", name))
	}

	ni := &orderedInitializer{
		order:          order,
		SubInitializer: i,
	}

	initializers = append(initializers, ni)

	cache[name] = ni
}
