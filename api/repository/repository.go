package repository

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewUserRoleRepository),
	fx.Provide(NewRoleRepository),
	fx.Provide(NewRoleMenuRepository),
	fx.Provide(NewMenuRepository),
	fx.Provide(NewMenuActionRepository),
	fx.Provide(NewMenuActionResourceRepository),
)
