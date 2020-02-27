package app

import (
	"net/http"

	"github.com/tjsampson/token-svc/internal/httphelper"
	"github.com/tjsampson/token-svc/internal/serviceprovider"
)

func listUsersHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	appCtxProvider.Logger.For(req.Context()).Info("entering listUsersHandler")

	users, err := appCtxProvider.UserService.List(req.Context())
	if err != nil {
		return httphelper.AppErr(err, "listUsersHandler.UserService.List")
	}

	appCtxProvider.Logger.For(req.Context()).Info("leaving listUsersHandler")
	return httphelper.AppResponse(http.StatusOK, users)
}
