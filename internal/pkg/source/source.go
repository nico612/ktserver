package source

import (
	"errors"
	"ktserver/internal/pkg/initdb"
)

const (
	initOrderSystem = 10

	initOrderApi        = initOrderSystem + 1
	initOrderCasbin     = initOrderApi + 1
	initOrderAuthority  = initOrderCasbin + 1
	initOrderDict       = initOrderCasbin + 1
	initOrderMenu       = initOrderAuthority + 1
	initOrderDictDetail = initOrderDict + 1
	initOrderUser       = initOrderAuthority + 1

	initOrderExcelTemplate = initOrderDictDetail + 1.
	initOrderMenuAuthority = initOrderMenu + initOrderAuthority
)

var (
	ErrMissingDBContext        = errors.New("missing db in context")
	ErrMissingDependentContext = errors.New("missing dependent value in context")
)

func RegisterSource() {
	// 在开发后台管理系统时，也应该按照如下顺序创建表
	initdb.RegisterInit(initOrderApi, &InitApi{})
	initdb.RegisterInit(initOrderCasbin, &initCasbin{})
	initdb.RegisterInit(initOrderAuthority, &initAuthority{})
	initdb.RegisterInit(initOrderMenu, &initMenu{})
	initdb.RegisterInit(initOrderDict, &initDict{})
	initdb.RegisterInit(initOrderDictDetail, &initDictDetail{})
	initdb.RegisterInit(initOrderExcelTemplate, &initExcelTemplate{})
	initdb.RegisterInit(initOrderUser, &initUser{})
	initdb.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}
