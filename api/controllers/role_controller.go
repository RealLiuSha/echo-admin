package controllers

import (
	"net/http"

	"github.com/RealLiuSha/echo-admin/api/services"
	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models"
	"github.com/RealLiuSha/echo-admin/models/dto"
	"github.com/RealLiuSha/echo-admin/pkg/echox"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type RoleController struct {
	logger      lib.Logger
	roleService services.RoleService
}

// NewRoleController creates new role controller
func NewRoleController(
	logger lib.Logger,
	roleService services.RoleService,
) RoleController {
	return RoleController{
		logger:      logger,
		roleService: roleService,
	}
}

// @tags Role
// @summary Role Query
// @produce application/json
// @param data query models.RoleQueryParam true "RoleQueryParam"
// @success 200 {object} echox.Response{data=models.RoleQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles [get]
func (a RoleController) Query(ctx echo.Context) error {
	param := new(models.RoleQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.roleService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Role
// @summary Role Get All
// @produce application/json
// @param data query models.RoleQueryParam true "RoleQueryParam"
// @success 200 {object} echox.Response{data=models.Roles} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles [get]
func (a RoleController) GetAll(ctx echo.Context) error {
	qr, err := a.roleService.Query(&models.RoleQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// @tags Role
// @summary Role Get By ID
// @produce application/json
// @param id path int true "role id"
// @success 200 {object} echox.Response{data=models.Role} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles/{id} [get]
func (a RoleController) Get(ctx echo.Context) error {
	role, err := a.roleService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: role}.JSON(ctx)
}

// @tags Role
// @summary Role Create
// @produce application/json
// @param data body models.Role true "Role"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles [post]
func (a RoleController) Create(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	role.CreatedBy = claims.Username

	id, err := a.roleService.WithTrx(trxHandle).Create(role)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Role
// @summary Role Update By ID
// @produce application/json
// @param id path int true "role id"
// @param data body models.Role true "Role"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles/{id} [put]
func (a RoleController) Update(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.roleService.WithTrx(trxHandle).Update(ctx.Param("id"), role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Role
// @summary Role Delete By ID
// @produce application/json
// @param id path int true "role id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles/{id} [delete]
func (a RoleController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.roleService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Role
// @summary Role Enable By ID
// @produce application/json
// @param id path int true "role id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles/{id}/enable [patch]
func (a RoleController) Enable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Role
// @summary Role Disable By ID
// @produce application/json
// @param id path int true "role id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/roles/{id}/disable [patch]
func (a RoleController) Disable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
