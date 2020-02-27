package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/middleware"
	"github.com/tjsampson/token-svc/internal/serviceprovider"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// App is the interface for the Application
type App interface {
	Serve() error
}

// app is the main server that serves up the API
type app struct {
	done       chan bool
	sigChannel chan os.Signal
	healthy    int32
	router     *mux.Router
	server     *http.Server
	logger     log.Factory
	appCtx     *serviceprovider.Context
}

// New creates a new Application server (our Rest API)
// this is where the mux router is created and decorated with middleware
func New(appCtxProvider *serviceprovider.Context) App {
	router := mux.NewRouter().StrictSlash(true)

	a := &app{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", appCtxProvider.Config.API.Port),
			Handler: handlers.CORS(
				handlers.AllowedOrigins(appCtxProvider.Config.API.AllowedOrigins),
				handlers.AllowedHeaders(appCtxProvider.Config.API.AllowedHeaders),
				handlers.AllowedMethods(appCtxProvider.Config.API.AllowedMethods),
			)(middleware.Adapt(
				router,
				middleware.AuthHandler(appCtxProvider),
				middleware.LogMetricsHandler(appCtxProvider.Logger, appCtxProvider.Metrics),
				middleware.TimeoutHandler(appCtxProvider.Config.API.TimeoutSecs),
				middleware.TracingHandler(appCtxProvider))),
			ReadTimeout:  time.Duration(appCtxProvider.Config.API.ReadTimeOutSecs) * time.Second,
			WriteTimeout: time.Duration(appCtxProvider.Config.API.WriteTimeOutSecs) * time.Second,
			IdleTimeout:  time.Duration(appCtxProvider.Config.API.IdleTimeOutSecs) * time.Second,
		},
		done:       make(chan bool),
		logger:     appCtxProvider.Logger.With(zap.String("package", "server")),
		sigChannel: make(chan os.Signal, 1024),
		healthy:    0,
		router:     router,
		appCtx:     appCtxProvider,
	}

	a.registerRoutes(appCtxProvider)
	return a
}

// Handle SIGNALS
func (a *app) sigHandler() {
	for {
		sig := <-a.sigChannel
		switch sig {
		case syscall.SIGHUP:
			a.logger.Bg().Info("reload config not setup up")
		case os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
			a.logger.Bg().Info("caught shutdown signal", zap.String("signal", sig.String()))
			a.gracefulShutdown()
		}
	}
}

func (a *app) gracefulShutdown() {
	// Store atomic "health" value
	// Prevents new requests from coming in while draining
	atomic.StoreInt32(&a.healthy, 0)

	// Pause the Context for `ShutdownTimeoutSecs` config value
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.appCtx.Config.API.ShutdownTimeoutSecs)*time.Second)
	defer cancel()

	// Turn off keepalive
	a.server.SetKeepAlivesEnabled(false)

	// Attempt to shutdown cleanly
	if err := a.server.Shutdown(ctx); err != nil {
		// YIKES! Shutdown failed, time to panic
		panic("http server failed graceful shutdown")
	}
	close(a.done)
}

// Serve is the entrypoint into the api
func (a *app) Serve() error {
	// Serve Up Metrics on a separte port
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%v", a.appCtx.Config.API.MetricsPort), nil)

	// signal the token-svc channel whenever an OS.Interrupt or SIGHUP occur
	// (both currently terminate. would like to use the SIGHUP for config reload)
	signal.Notify(a.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// goroutine to handle signals
	go a.sigHandler()

	// atomically store the health as "healthy=1"
	atomic.StoreInt32(&a.healthy, 1)

	// log server start details
	a.logger.Bg().Info("server up ===========>",
		zap.String("port", a.appCtx.Config.API.Port),
		zap.String("metrics_port", a.appCtx.Config.API.MetricsPort),
		zap.String("full_version", a.appCtx.VersionInfo.FullVersionNumber(true)),
		zap.String("build_date", a.appCtx.VersionInfo.BuildDate),
		zap.String("metadata", a.appCtx.VersionInfo.VersionMetadata),
		zap.String("prerelease", a.appCtx.VersionInfo.VersionPrerelease),
		zap.String("version", a.appCtx.VersionInfo.Version),
		zap.String("revision", a.appCtx.VersionInfo.Revision),
		zap.String("author", a.appCtx.VersionInfo.Author),
		zap.String("branch", a.appCtx.VersionInfo.Branch),
		zap.String("builder", a.appCtx.VersionInfo.BuildUser),
		zap.String("host", a.appCtx.VersionInfo.BuildHost))

	// Tally Build Metrics
	a.appCtx.Metrics.StatBuildInfo.WithLabelValues(
		a.appCtx.Config.API.ServiceName,
		a.appCtx.VersionInfo.Revision,
		a.appCtx.VersionInfo.Branch,
		a.appCtx.VersionInfo.Version,
		a.appCtx.VersionInfo.Author,
		a.appCtx.VersionInfo.BuildDate,
		a.appCtx.VersionInfo.BuildUser,
		a.appCtx.VersionInfo.BuildHost).Set(1)

	// serve up the api by listening on the configured port
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Bg().Error("failed ListenAndServe", zap.String("port", a.appCtx.Config.API.Port))
		return err
	}

	// signal caught in a.gracefulShutdown() and close(a.done) called
	<-a.done

	// log server shutdown details
	a.logger.Bg().Info("graceful server shutdown  ===========>", zap.String("lifetime", a.appCtx.VersionInfo.UpTime().String()))
	return nil
}
