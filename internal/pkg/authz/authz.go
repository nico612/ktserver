package authz

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

// CasbinAuthorizer 授权
type CasbinAuthorizer struct {
	Enforcer *casbin.SyncedCachedEnforcer
	log      *log.Helper
}

var syncedCachedEnforcer *casbin.SyncedCachedEnforcer

// NewCasbinAuthorizer 实例化CasbinAuthorizer
func NewCasbinAuthorizer(db *gorm.DB, logger log.Logger) *CasbinAuthorizer {
	if syncedCachedEnforcer == nil {
		once.Do(func() {
			syncedCachedEnforcer = newEnforcer(db)
		})
	}
	return &CasbinAuthorizer{
		Enforcer: syncedCachedEnforcer,
		log:      log.NewHelper(logger),
	}
}

// newEnforcer 初始化casbin
func newEnforcer(db *gorm.DB) *casbin.SyncedCachedEnforcer {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}

	// 模型配置
	// request_definition: 定义请求的数据结构 r = sub, obj, act, sub: 用户, obj: 资源, act: 操作
	// policy_definition: 定义策略的数据结构 p = sub, obj, act, sub: 用户, obj: 资源, act: 操作
	// role_definition: 定义角色的数据结构 g = _, _, g: 用户, _: 角色
	// policy_effect: 策略结果 e = some(where (p.eft == allow)), e: 策略结果
	// matchers: 匹配规则 m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	// 加载模型
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(err)
	}

	// 创建casbin
	enforcer, err := casbin.NewSyncedCachedEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	// 设置缓存过期时间
	enforcer.SetExpireTime(60 * 60)

	// 加载策略
	_ = enforcer.LoadPolicy()

	return enforcer
}
