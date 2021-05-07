package controllers

import (
	"net/http"
	"strings"

	"github.com/RealLiuSha/echo-admin/api/services"
	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/models/dto"
	"github.com/RealLiuSha/echo-admin/pkg/echox"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type UserController struct {
	userService services.UserService
	logger      lib.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger lib.Logger) UserController {
	return UserController{
		userService: userService,
		logger:      logger,
	}
}

// @tags User
// @summary User Query
// @produce application/json
// @param data query models.UserQueryParam true "UserQueryParam"
// @success 200 {object} echox.Response{data=models.UserQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users [get]
func (a UserController) Query(ctx echo.Context) error {
	param := new(models.UserQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}
	if v := ctx.QueryParam("role_ids"); v != "" {
		param.RoleIDs = strings.Split(v, ",")
	}

	qr, err := a.userService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags User
// @summary User Create
// @produce application/json
// @param data body models.User true "User"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users [post]
func (a UserController) Create(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	} else if user.Password == "" {
		return echox.Response{Code: http.StatusBadRequest, Message: errors.UserPasswordRequired}.JSON(ctx)
	}

	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	user.CreatedBy = claims.Username

	qr, err := a.userService.WithTrx(trxHandle).Create(user)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags User
// @summary User Get By ID
// @produce application/json
// @param id path int true "user id"
// @success 200 {object} echox.Response{data=models.User} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users/{id} [get]
func (a UserController) Get(ctx echo.Context) error {
	user, err := a.userService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: user}.JSON(ctx)
}

// @tags User
// @summary User Update By ID
// @produce application/json
// @param id path int true "user id"
// @param data body models.User true "User"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users/{id} [put]
func (a UserController) Update(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	err := a.userService.WithTrx(trxHandle).Update(ctx.Param("id"), user)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags User
// @summary User Delete By ID
// @produce application/json
// @param id path int true "user id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users/{id} [delete]
func (a UserController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	err := a.userService.WithTrx(trxHandle).Delete(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags User
// @summary User Enable By ID
// @produce application/json
// @param id path int true "user id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users/{id}/enable [patch]
func (a UserController) Enable(ctx echo.Context) error {
	err := a.userService.UpdateStatus(ctx.Param("id"), 1)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags User
// @summary User Disable By ID
// @produce application/json
// @param id path int true "user id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/users/{id}/disable [patch]
func (a UserController) Disable(ctx echo.Context) error {
	err := a.userService.UpdateStatus(ctx.Param("id"), -1)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
