package healthservice

import (
	"context"
	"runtime"
	"time"

	"github.com/tjsampson/token-svc/internal/datastores/redis"
	"github.com/tjsampson/token-svc/internal/errors"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/healthmodels"
	"github.com/tjsampson/token-svc/internal/repos/healthrepo"
	"github.com/tjsampson/token-svc/pkg/helpers/datahelpers"
	"github.com/tjsampson/token-svc/pkg/metrics"
	"github.com/tjsampson/token-svc/pkg/version"

	"go.uber.org/zap"
)

// Service is the health service Interface
type Service interface {
	GetMemoryStats(ctx context.Context) (healthmodels.MemoryStats, error)
	GetCacheHealth(ctx context.Context) (healthmodels.CacheHealth, error)
	GetDatabaseHealth(ctx context.Context) (healthmodels.DBHealth, error)
	GetAPIHealth(ctx context.Context) (healthmodels.APIHealth, error)
	GetAPIUpTime(ctx context.Context) healthmodels.APIUpTime
	GetFullHealth(ctx context.Context) (healthmodels.TokenSvcHealth, error)
}

// New returns a new healthRepo store
func New(
	healthRepo healthrepo.Store,
	redisClient redis.Provider,
	versionInfo version.Info,
	metricProvider *metrics.Provider,
	logger log.Factory) Service {

	return &service{
		logger:         logger.With(zap.String("package", "healthservice")),
		repo:           healthRepo,
		redisClient:    redisClient,
		versionInfo:    versionInfo,
		metricProvider: metricProvider,
	}

}

type service struct {
	logger         log.Factory
	repo           healthrepo.Store
	redisClient    redis.Provider
	versionInfo    version.Info
	metricProvider *metrics.Provider
}

func (svc *service) GetAPIUpTime(ctx context.Context) healthmodels.APIUpTime {
	uptimeSec := time.Since(svc.versionInfo.ServerStartTime).Seconds()
	uptimeMins := uptimeSec / 60
	uptimeHours := uptimeMins / 60
	uptimeDays := uptimeHours / 24
	uptimeWeeks := uptimeDays / 7
	return healthmodels.APIUpTime{
		Seconds: uptimeSec,
		Minutes: uptimeMins,
		Hours:   uptimeHours,
		Days:    uptimeDays,
		Weeks:   uptimeWeeks,
	}
}

func (svc *service) GetMemoryStats(ctx context.Context) (healthmodels.MemoryStats, error) {
	memStats := memStats()
	svc.metricProvider.StatMemAllocGuage.Set(float64(memStats.Allocation))
	svc.metricProvider.StatMemTotalAllocGuage.Set(float64(memStats.TotalAllocation))
	svc.metricProvider.StatMemSysGuage.Set(float64(memStats.Sys))
	svc.metricProvider.StatMemNumGCGuage.Set(float64(memStats.NumGC))
	svc.metricProvider.StatGoRoutineGuage.Set(float64(runtime.NumGoroutine()))
	return memStats, nil
}

func (svc *service) GetCacheHealth(ctx context.Context) (healthmodels.CacheHealth, error) {
	var err error
	var result healthmodels.CacheHealth

	// ping the redis cache to check it's health
	result.Ping, err = svc.redisClient.Ping(ctx)

	return result, errors.ErrorWrapper(err, "HealthService.GetCacheHealth")
}

func (svc *service) GetDatabaseHealth(ctx context.Context) (healthmodels.DBHealth, error) {
	var err error
	var result healthmodels.DBHealth

	// Dip the db to check it's health
	result, err = svc.repo.HealthStatus(ctx)

	return result, errors.ErrorWrapper(err, "HealthService.GetDatabaseHealth")
}

func (svc *service) GetAPIHealth(ctx context.Context) (healthmodels.APIHealth, error) {
	return svc.apiHealth(ctx), nil
}

func (svc *service) GetFullHealth(ctx context.Context) (healthmodels.TokenSvcHealth, error) {
	var err error
	var result healthmodels.TokenSvcHealth

	cacheHealth, err := svc.GetCacheHealth(ctx)
	if err != nil {
		return result, err
	}

	dbHealth, err := svc.repo.HealthStatus(ctx)
	if err != nil {
		return result, errors.ErrorWrapper(err, "HealthService.GetFullHealth.repo.healthStatus")
	}

	result = healthmodels.TokenSvcHealth{
		APIHealth:   svc.apiHealth(ctx),
		DBHealth:    dbHealth,
		MemoryStats: memStats(),
		CacheHealth: cacheHealth,
	}

	return result, errors.ErrorWrapper(err, "HealthService.GetFullHealth")
}

func (svc *service) apiHealth(ctx context.Context) healthmodels.APIHealth {
	return healthmodels.APIHealth{
		Status:  "OK",
		Version: svc.versionInfo,
		Uptime:  svc.GetAPIUpTime(ctx),
	}
}

func memStats() healthmodels.MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return healthmodels.MemoryStats{
		Allocation:      datahelpers.BytesToMb(m.Alloc),
		TotalAllocation: datahelpers.BytesToMb(m.TotalAlloc),
		Sys:             datahelpers.BytesToMb(m.Sys),
		NumGC:           m.NumGC,
		GoRoutines:      runtime.NumGoroutine(),
	}
}
