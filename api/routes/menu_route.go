package routes

import (
	"github.com/RealLiuSha/echo-admin/api/controllers"
	"github.com/RealLiuSha/echo-admin/lib"
)

type MenuRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	menuController controllers.MenuController
}

// NewMenuRoutes creates new menu routes
func NewMenuRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	menuController controllers.MenuController,
) MenuRoutes {
	return MenuRoutes{
		handler:        handler,
		logger:         logger,
		menuController: menuController,
	}
}

// Setup menu routes
func (a MenuRoutes) Setup() {
	a.logger.Zap.Info("Setting up menu routes")
	api := a.handler.RouterV1.Group("/menus")
	{
		api.GET("", a.menuController.Query)

		api.POST("", a.menuController.Create)
		api.GET("/:id", a.menuController.Get)
		api.PUT("/:id", a.menuController.Update)
		api.DELETE("/:id", a.menuController.Delete)
		api.PATCH("/:id/enable", a.menuController.Enable)
		api.PATCH("/:id/disable", a.menuController.Disable)

		api.GET("/:id/actions", a.menuController.GetActions)
		api.PUT("/:id/actions", a.menuController.UpdateActions)
	}
}
