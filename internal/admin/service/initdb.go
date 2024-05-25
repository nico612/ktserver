package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"ktserver/internal/admin/conf"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/initdb"
	"ktserver/internal/pkg/response"
)

// DBUService 初始化数据库
type DBUService struct {
	c             *conf.Config
	log           *log.Helper
	initdbService *initdb.InitDBService
}

// NewDBUService  初始化数据库
func NewDBUService(c *conf.Config, logger log.Logger, initdbService *initdb.InitDBService) *DBUService {
	return &DBUService{
		c:             c,
		log:           log.NewHelper(log.With(logger, "service", "DBUService")),
		initdbService: initdbService,
	}
}

// CheckDB 检测是否需要初始化数据库
func (uc *DBUService) CheckDB(c *gin.Context) {

}

// InitDB
// @Tags     InitDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Param    data  body      initdb.InitDBModel                  true  "初始化数据库参数"
// @Success  200   {object}  response.Response{data=string}  "初始化用户数据库"
// @Router   /init/initdb [post]
func (uc *DBUService) InitDB(c *gin.Context) {
	var dbInfo initdb.InitDBModel
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		log.Errorf("参数校验不通过! error = %+v", err)
		response.Result(c, bizerr.InvalidParam.Wrap(err))
		return
	}
	//
	//initService := initdb.NewInitDBService(
	//	handler.NewMysqlInitHandler(),
	//)
	if err := uc.initdbService.InitDB(&dbInfo); err != nil {
		uc.log.Errorf("自动创建数据库失败! error = %+v", err)
		response.Result(c, err)
		return
	}

	// 创建数据库成功
	response.Result(c, "初始化数据库成功")
}
