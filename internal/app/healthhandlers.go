package app

import (
	"net/http"

	"github.com/tjsampson/token-svc/internal/httphelper"
	"github.com/tjsampson/token-svc/internal/serviceprovider"
)

func pingHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	return httphelper.AppResponse(http.StatusOK, "pong")
}

func apiHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	result, err := appCtxProvider.HealthService.GetAPIHealth(req.Context())

	if err != nil {
		return httphelper.AppErr(err, "apiHealthHandler")
	}

	return httphelper.AppResponse(http.StatusOK, result)
}

func getFullHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {

	result, err := appCtxProvider.HealthService.GetFullHealth(req.Context())

	if err != nil {
		return httphelper.AppErr(err, "getFullHealthHandler")
	}

	return httphelper.AppResponse(http.StatusOK, result)
}

func databaseHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	result, err := appCtxProvider.HealthService.GetDatabaseHealth(req.Context())

	if err != nil {
		return httphelper.AppErr(err, "databaseHealthHandler")
	}

	return httphelper.AppResponse(http.StatusOK, result)
}

func memoryHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	result, err := appCtxProvider.HealthService.GetMemoryStats(req.Context())

	if err != nil {
		return httphelper.AppErr(err, "memoryHealthHandler")
	}

	return httphelper.AppResponse(http.StatusOK, result)
}

func cacheHealthHandler(appCtxProvider *serviceprovider.Context, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	result, err := appCtxProvider.HealthService.GetCacheHealth(req.Context())

	if err != nil {
		return httphelper.AppErr(err, "cacheHealthHandler")
	}

	return httphelper.AppResponse(http.StatusOK, result)
}
