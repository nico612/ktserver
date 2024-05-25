package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewBaseUseCase,      // 基础
	NewCaptchaUseCase,   // 验证码
	NewUserUseCase,      // 用户
	NewMenuUseCase,      // 菜单
	NewAuthorityUseCase, // 角色
)
