package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewPprofRoutes),
	fx.Provide(NewSwaggerRoutes),
	fx.Provide(NewPublicRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewRoleRoutes),
	fx.Provide(NewMenuRoutes),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	pprofRoutes PprofRoutes,
	swaggerRoutes SwaggerRoutes,
	publicRoutes PublicRoutes,
	userRoutes UserRoutes,
	roleRoutes RoleRoutes,
	menuRoutes MenuRoutes,
) Routes {
	return Routes{
		pprofRoutes,
		swaggerRoutes,
		publicRoutes,
		userRoutes,
		roleRoutes,
		menuRoutes,
	}
}

// Setup all the route
func (a Routes) Setup() {
	for _, route := range a {
		route.Setup()
	}
}
