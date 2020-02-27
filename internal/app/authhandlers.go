package app

import (
	"net/http"

	"github.com/tjsampson/token-svc/internal/httphelper"
	"github.com/tjsampson/token-svc/internal/models/authmodels"
	"github.com/tjsampson/token-svc/internal/serviceprovider"

	"go.uber.org/zap"
)

func loginHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	appCtxProvider.Logger.For(req.Context()).Info("entering loginHandler")
	userCreds := &authmodels.UserCreds{}

	if err := httphelper.ParseBody(res, req, userCreds); err != nil {
		return httphelper.AppErr(err, "loginHandler.ParseBody")
	}

	if err := appCtxProvider.Validator.Validate(userCreds); err != nil {
		return httphelper.AppErr(err, "loginHandler.Validate")
	}

	loginResults, err := appCtxProvider.AuthService.Login(req.Context(), userCreds)

	if err != nil {
		return httphelper.AppErr(err, "loginHandler.AuthService.Login")
	}

	http.SetCookie(res, loginResults.HTTPCookie)
	appCtxProvider.Logger.For(req.Context()).Info("leaving loginHandler", zap.String("email", userCreds.Email))
	return httphelper.AppResponse(http.StatusOK, map[string]string{"access_token": loginResults.AccessToken, "refresh_token": loginResults.RefreshToken})
}

func registerHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	appCtxProvider.Logger.For(req.Context()).Info("entering registerHandler")

	userCreds := &authmodels.UserRegistration{}
	var err error

	if err = httphelper.ParseBody(res, req, userCreds); err != nil {
		return httphelper.AppErr(err, "registerHandler.ParseBody")
	}

	if err = appCtxProvider.Validator.Validate(userCreds); err != nil {
		return httphelper.AppErr(err, "registerHandler.Validate")
	}

	user, err := appCtxProvider.AuthService.Register(req.Context(), userCreds)
	if err != nil {
		return httphelper.AppErr(err, "registerHandler.AuthService.Register")
	}

	appCtxProvider.Logger.For(req.Context()).Info("leaving registerHandler", zap.String("email", userCreds.Email))
	return httphelper.AppResponse(http.StatusCreated, user)
}
