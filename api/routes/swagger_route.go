package routes

import (
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/docs"
	"github.com/RealLiuSha/echo-admin/lib"
)

// @securityDefinitions.apikey X-Auth-Token
// @in header
// @name Authorization
// @schemes http https
// @basePath /
// @contact.name LiuSha
// @contact.email liusha@email.cn
type SwaggerRoutes struct {
	config  lib.Config
	logger  lib.Logger
	handler lib.HttpHandler
}

// NewUserRoutes creates new swagger routes
func NewSwaggerRoutes(
	config lib.Config,
	logger lib.Logger,
	handler lib.HttpHandler,
) SwaggerRoutes {
	return SwaggerRoutes{
		config:  config,
		logger:  logger,
		handler: handler,
	}
}

// Setup swagger routes
func (a SwaggerRoutes) Setup() {
	docs.SwaggerInfo.Title = a.config.Name
	docs.SwaggerInfo.Version = constants.Version

	a.logger.Zap.Info("Setting up swagger routes")
	a.handler.Engine.GET("/swagger/*", echoSwagger.WrapHandler)
}
