package routes

import (
	"github.com/RealLiuSha/echo-admin/api/controllers"
	"github.com/RealLiuSha/echo-admin/lib"
)

type PublicRoutes struct {
	logger            lib.Logger
	handler           lib.HttpHandler
	publicController  controllers.PublicController
	captchaController controllers.CaptchaController
}

// NewUserRoutes creates new public routes
func NewPublicRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	publicController controllers.PublicController,
	captchaController controllers.CaptchaController,
) PublicRoutes {
	return PublicRoutes{
		handler:           handler,
		logger:            logger,
		publicController:  publicController,
		captchaController: captchaController,
	}
}

// Setup public routes
func (a PublicRoutes) Setup() {
	a.logger.Zap.Info("Setting up public routes")
	api := a.handler.RouterV1.Group("/publics")
	{
		api.GET("/user", a.publicController.UserInfo)
		api.POST("/user/login", a.publicController.UserLogin)
		api.POST("/user/logout", a.publicController.UserLogout)
		api.GET("/user/menutree", a.publicController.MenuTree)
		//api.GET("/user/password", a.publicController.UserPassword)

		// sys routes
		api.GET("/sys/routes", a.publicController.SysRoutes)

		// captcha
		api.GET("/captcha", a.captchaController.GetCaptcha)
		api.POST("/captcha/verify", a.captchaController.VerifyCaptcha)
	}
}
