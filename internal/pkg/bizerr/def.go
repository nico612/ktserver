package bizerr

var (
	Success      = New(0, "success")
	UnknownError = New(7, "internal error")
	InvalidParam = New(100001, "invalid parameter")

	// BaseBizErr
	CaptchaLimit       = New(100002, "图形验证码获取太频繁，请稍后再试")
	CaptchaInvalid     = New(100003, "图形验证码错误")
	CaptchaExpired     = New(100004, "图形验证码已过期")
	CaptchaBlocked     = New(100005, "图形验证码已被锁定，请稍后再试")
	LoginPasswordError = New(100006, "密码错误")

	ErrTokenInvalid         = New(100007, "无效的token")
	ErrTokenExpired         = New(100008, "token已过期")
	ErrTokenParseFail       = New(100009, "解析token失败")
	ErrSignTokenFailed      = New(100011, "签名token失败")
	UnAuthority             = New(100012, "未登录或非法访问")
	ErrUserLoginOtherDevice = New(100013, "当前账号已在其他设备登录")
	ErrNoPermission         = New(100014, "权限不足")

	// user
	InvalidPassword = New(100101, "密码错误")
	UserNotFound    = New(100102, "用户不存在")

	// 菜单
	MenuNotFound = New(100201, "菜单不存在")
	MenuExist    = New(100202, "菜单已存在")
)
