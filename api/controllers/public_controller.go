package controllers

import (
	"net/http"
	"strings"

	"github.com/RealLiuSha/echo-admin/api/services"
	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/models/dto"
	"github.com/RealLiuSha/echo-admin/pkg/echox"
	"github.com/labstack/echo/v4"
)

type PublicController struct {
	userService services.UserService
	authService services.AuthService
	captcha     lib.Captcha
	logger      lib.Logger
}

// NewPublicController creates new public controller
func NewPublicController(
	userService services.UserService,
	authService services.AuthService,
	captcha lib.Captcha,
	logger lib.Logger,
) PublicController {
	return PublicController{
		userService: userService,
		authService: authService,
		captcha:     captcha,
		logger:      logger,
	}
}

type route struct {
	*echo.Route
	Name *struct{} `json:"name,omitempty"`
}

// @Tags Public
// @Summary UserLogin
// @Produce application/json
// @Param data body dto.Login true "Login"
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/publics/sys/routes [post]
func (a PublicController) SysRoutes(ctx echo.Context) error {
	routes := make([]*route, 0)
	for _, eRoute := range ctx.Echo().Routes() {
		// Only interfaces starting with /api/ are exposed
		if !strings.HasPrefix(eRoute.Path, "/api/") {
			continue
		}
		routes = append(routes, &route{Route: eRoute})
	}

	return echox.Response{Code: http.StatusOK, Data: routes}.JSON(ctx)
}

// @Tags Public
// @Summary UserInfo
// @Produce application/json
// @Success 200 {string} echox.Response{data=models.UserInfo} "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/publics/user [get]
func (a PublicController) UserInfo(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	userinfo, err := a.userService.GetUserInfo(claims.ID)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: userinfo}.JSON(ctx)
}

// @Tags Public
// @Summary UserMenuTree
// @Produce application/json
// @Success 200 {string} echox.Response{data=models.MenuTrees} "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/publics/user/menutree [get]
func (a PublicController) MenuTree(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	menuTrees, err := a.userService.GetUserMenuTrees(claims.ID)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: menuTrees}.JSON(ctx)
}

// @Tags Public
// @Summary UserLogin
// @Produce application/json
// @Param data body dto.Login true "Login"
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/publics/user/login [post]
func (a PublicController) UserLogin(ctx echo.Context) error {
	login := new(dto.Login)

	if err := ctx.Bind(login); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if !a.captcha.Verify(login.CaptchaID, login.CaptchaCode, false) {
		return echox.Response{Code: http.StatusBadRequest, Message: errors.CaptchaAnswerCodeNoMatch}.JSON(ctx)
	}

	user, err := a.userService.Verify(login.Username, login.Password)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	token, err := a.authService.GenerateToken(user)
	if err != nil {
		return echox.Response{Code: http.StatusInternalServerError, Message: errors.AuthTokenGenerateFail}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"token": token}}.JSON(ctx)
}

// @Tags Public
// @Summary UserLogout
// @Produce application/json
// @Success 200 {string} echox.Response "success"
// @Router /api/publics/user/logout [post]
func (a PublicController) UserLogout(ctx echo.Context) error {
	claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	if ok {
		a.authService.DestroyToken(claims.Username)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}
