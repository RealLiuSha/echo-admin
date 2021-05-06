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

type MenuController struct {
	menuService services.MenuService
	logger      lib.Logger
}

// NewMenuController creates new menu controller
func NewMenuController(
	logger lib.Logger,
	menuService services.MenuService,
) MenuController {
	return MenuController{
		logger:      logger,
		menuService: menuService,
	}
}

// @tags Menu
// @summary Menu Query
// @produce application/json
// @param data query models.MenuQueryParam true "MenuQueryParam"
// @success 200 {object} echox.Response{data=models.MenuQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus [get]
func (a MenuController) Query(ctx echo.Context) error {
	param := new(models.MenuQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := a.menuService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if param.Tree {
		return echox.Response{Code: http.StatusOK, Data: qr.List.ToMenuTrees()}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// @tags Menu
// @summary Menu Get By ID
// @produce application/json
// @param id path int true "menu id"
// @success 200 {object} echox.Response{data=models.Menu} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id} [get]
func (a MenuController) Get(ctx echo.Context) error {
	menu, err := a.menuService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: menu}.JSON(ctx)
}

// @tags Menu
// @summary Menu Create
// @produce application/json
// @param data body models.Menu true "Menu"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus [post]
func (a MenuController) Create(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	menu.CreatedBy = claims.Username

	id, err := a.menuService.WithTrx(trxHandle).Create(menu)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// @tags Menu
// @summary Menu Update By ID
// @produce application/json
// @param id path int true "menu id"
// @param data body models.Menu true "Menu"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id} [put]
func (a MenuController) Update(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).Update(ctx.Param("id"), menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Menu
// @summary Menu Delete By ID
// @produce application/json
// @param id path int true "menu id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id} [delete]
func (a MenuController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Menu
// @summary Menu Enable By ID
// @produce application/json
// @param id path int true "menu id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id}/enable [patch]
func (a MenuController) Enable(ctx echo.Context) error {
	if err := a.menuService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Menu
// @summary Menu Disable By ID
// @produce application/json
// @param id path int true "menu id"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id}/disable [patch]
func (a MenuController) Disable(ctx echo.Context) error {
	if err := a.menuService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// @tags Menu
// @summary MenuActions Get By menuID
// @produce application/json
// @param id path int true "menu id"
// @success 200 {object} echox.Response{data=models.MenuActions} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id}/actions [get]
func (a MenuController) GetActions(ctx echo.Context) error {
	actions, err := a.menuService.GetMenuActions(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: actions}.JSON(ctx)
}

// @tags Menu
// @summary Menu Actions Update By menuID
// @produce application/json
// @param id path int true "menu id"
// @param data body models.MenuActions true "Menu"
// @success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @router /api/menus/{id}/actions [put]
func (a MenuController) UpdateActions(ctx echo.Context) error {
	actions := make(models.MenuActions, 0)
	if err := ctx.Bind(&actions); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := a.menuService.WithTrx(trxHandle).UpdateActions(ctx.Param("id"), actions); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
