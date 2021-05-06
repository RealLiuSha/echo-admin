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

func (a RoleController) GetAll(ctx echo.Context) error {
	qr, err := a.roleService.Query(&models.RoleQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})

	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

func (a RoleController) Get(ctx echo.Context) error {
	role, err := a.roleService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: role}.JSON(ctx)
}

func (a RoleController) Create(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(dto.JwtClaims)
	role.CreatedBy = claims.Username

	id, err := a.roleService.WithTrx(trxHandle).Create(role)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

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

func (a RoleController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.roleService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

func (a RoleController) Enable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

func (a RoleController) Disable(ctx echo.Context) error {
	if err := a.roleService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
