package controllers

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewPublicController),
	fx.Provide(NewCaptchaController),
	fx.Provide(NewUserController),
	fx.Provide(NewRoleController),
	fx.Provide(NewMenuController),
)
