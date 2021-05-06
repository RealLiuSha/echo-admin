package services

import "go.uber.org/fx"

// Module exports services present
var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewRoleService),
	fx.Provide(NewMenuService),
	fx.Provide(NewCasbinService),
	fx.Provide(NewAuthService),
)
