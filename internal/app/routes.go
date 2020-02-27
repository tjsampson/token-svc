package app

import (
	"github.com/tjsampson/token-svc/internal/middleware"
	"github.com/tjsampson/token-svc/internal/serviceprovider"
)

func (a *app) registerRoutes(appCtxProvider *serviceprovider.Context) {
	a.router.Handle("/login", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: loginHandler}).Methods("POST")
	a.router.Handle("/register", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: registerHandler}).Methods("POST")
	a.router.Handle("/health", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: getFullHealthHandler}).Methods("GET")
	a.router.Handle("/health/api", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: apiHealthHandler}).Methods("GET")
	a.router.Handle("/health/ping", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: pingHealthHandler}).Methods("GET")
	a.router.Handle("/health/database", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: databaseHealthHandler}).Methods("GET")
	a.router.Handle("/health/cache", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: cacheHealthHandler}).Methods("GET")
	a.router.Handle("/health/memory", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: memoryHealthHandler}).Methods("GET")
	a.router.Handle("/users", middleware.Handler{AppCtx: appCtxProvider, RouteHandler: listUsersHandler}).Methods("GET")
}
