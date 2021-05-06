package middlewares

import (
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/labstack/echo/v4/middleware"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler lib.HttpHandler, logger lib.Logger) CorsMiddleware {
	return CorsMiddleware{
		handler: handler,
		logger:  logger,
	}
}

func (a CorsMiddleware) Setup() {
	a.logger.Zap.Info("Setting up cors middleware")

	a.handler.Engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
	}))
}
