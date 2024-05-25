package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"ktserver/internal/admin/biz"
	"ktserver/internal/admin/conf"
	"ktserver/internal/admin/model/params"
	"ktserver/internal/pkg/bizerr"
	"ktserver/internal/pkg/response"
)

type BaseService struct {
	c              *conf.Config
	log            *log.Helper
	baseUseCase    *biz.BaseUseCase
	captchaUseCase *biz.CaptchaUseCase
}

func NewBaseService(c *conf.Config, logger log.Logger, baseUseCase *biz.BaseUseCase, captchaUseCase *biz.CaptchaUseCase) *BaseService {
	return &BaseService{
		c:              c,
		log:            log.NewHelper(log.With(logger, "module", "service/base")),
		baseUseCase:    baseUseCase,
		captchaUseCase: captchaUseCase,
	}
}

// Login
// @Tags Base
// @Summary 登录
// @Produce application/json
// @Param data body params.Login true  “用户吗，密码， 验证码”"
// @Success  200   {object}  response.Response{data=params.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router /base/login [post]
func (uc *BaseService) Login(c *gin.Context) {
	var req params.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Result(c, bizerr.InvalidParam)
		return
	}

	// verify captcha
	if err := uc.captchaUseCase.VerifyCaptcha(c, req.CaptchaId, req.Captcha); err != nil {
		response.Result(c, err)
		return
	}

	resp, err := uc.baseUseCase.Login(c, req)
	if err != nil {
		uc.log.Errorf("Login error: %v", err)
		response.Result(c, err)
		return
	}
	response.Result(c, resp)
}

// TODO token刷新，刷新token不经常用，直接储存到数据库中即可，校验旧token，查询刷新token，然后创建新的token

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=params.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (uc *BaseService) Captcha(c *gin.Context) {
	ipAddr := c.ClientIP()
	captchaId, b64s, _, err := uc.captchaUseCase.GenerateCaptcha(c, ipAddr)
	if err != nil {
		uc.log.Errorf("GenerateCaptcha error: %v", err)
		response.Result(c, err)
		return
	}

	resp := params.SysCaptchaResponse{
		CaptchaId:     captchaId,
		PicPath:       b64s,
		CaptchaLength: uc.c.CaptchaOptions.KeyLength,
		OpenCaptcha:   true,
	}

	response.Result(c, resp)
}
