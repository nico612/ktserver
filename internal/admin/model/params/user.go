package params

import "ktserver/internal/admin/model"

// User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

type LoginResponse struct {
	User      model.SysUser `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expiresAt"`
}

type UserInfoResp struct {
	UserInfo model.SysUser `json:"userInfo"`
}

// Modify password structure
type ChangePasswordReq struct {
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}
