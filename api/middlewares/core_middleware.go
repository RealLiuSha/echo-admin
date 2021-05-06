package middlewares

import (
	"fmt"
	"runtime"

	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// core middleware is a functional extension to "echo",
// including database transactions and panic recovery
// and more
type CoreMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
	db      lib.Database
}

// statusInList function checks if context writer status is in provided list
func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// NewCoreMiddleware creates new database transactions middleware
func NewCoreMiddleware(handler lib.HttpHandler, logger lib.Logger, db lib.Database) CoreMiddleware {
	return CoreMiddleware{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}

func (a CoreMiddleware) core() echo.MiddlewareFunc {
	logger := a.logger.DesugarZap.With(zap.String("module", "core-mw"))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			txHandle := a.db.ORM.Begin()
			logger.Info("beginning database transaction")

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// recovery stack
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, false)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					logger.Error(msg)

					// rollback database transaction
					logger.Info("rolling back transaction due to panic")
					txHandle.Rollback()
					ctx.Error(err)
				}
			}()

			ctx.Set(constants.DBTransaction, txHandle)

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			code := ctx.Response().Status
			// rollback transaction on server errors
			if code >= 400 {
				msg := fmt.Sprintf("rolling back transaction due to status code: %d", code)
				logger.Info(msg)

				txHandle.Rollback()
			} else {
				a.logger.DesugarZap.Info("committing transactions")
				if err := txHandle.Commit().Error; err != nil {
					logger.Error(fmt.Sprintf("trx commit error: %v", err))
				}
			}

			return nil
		}
	}
}

func (a CoreMiddleware) Setup() {
	a.logger.Zap.Info("setting up core middleware")
	a.handler.Engine.Use(a.core())
}
