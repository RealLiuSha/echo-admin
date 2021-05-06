package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"

	"github.com/RealLiuSha/echo-admin/lib"
)

type PprofRoutes struct {
	logger  lib.Logger
	handler lib.HttpHandler
}

// NewUserRoutes creates new pprof routes
func NewPprofRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
) PprofRoutes {
	return PprofRoutes{
		handler: handler,
		logger:  logger,
	}
}

// Setup pprof routes
func (a PprofRoutes) Setup() {
	a.logger.Zap.Info("Setting up pprof routes")

	r := a.handler.Engine.Group("/pprof")
	{
		r.GET("/", handler(pprof.Index))
		r.GET("/allocs", handler(pprof.Handler("allocs").ServeHTTP))
		r.GET("/block", handler(pprof.Handler("block").ServeHTTP))
		r.GET("/cmdline", handler(pprof.Cmdline))
		r.GET("/goroutine", handler(pprof.Handler("goroutine").ServeHTTP))
		r.GET("/heap", handler(pprof.Handler("heap").ServeHTTP))
		r.GET("/mutex", handler(pprof.Handler("mutex").ServeHTTP))
		r.GET("/profile", handler(pprof.Profile))
		r.POST("/symbol", handler(pprof.Symbol))
		r.GET("/symbol", handler(pprof.Symbol))
		r.GET("/threadcreate", handler(pprof.Handler("threadcreate").ServeHTTP))
		r.GET("/trace", handler(pprof.Trace))
	}
}

func handler(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
