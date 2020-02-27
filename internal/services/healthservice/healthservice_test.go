package healthservice

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/tjsampson/token-svc/internal/datastores/redis"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/healthmodels"
	"github.com/tjsampson/token-svc/internal/repos/healthrepo"
	"github.com/tjsampson/token-svc/pkg/metrics"
	"github.com/tjsampson/token-svc/pkg/version"
)

type mockHealthRepo struct {
}

type mockRedisClient struct {
}

func (mrc *mockRedisClient) Close() error {
	return nil
}

func (mrc *mockRedisClient) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return nil
}

// 	Get(ctx context.Context, key string) (string, error)
func (mrc *mockRedisClient) Get(ctx context.Context, key string) (string, error) {
	return "id", nil
}

func (mrc *mockRedisClient) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (mhr *mockHealthRepo) HealthStatus(ctx context.Context) (healthmodels.DBHealth, error) {
	return healthmodels.DBHealth{}, nil
}

func MockHealthRepo() healthrepo.Store {
	return &mockHealthRepo{}
}

func MockRedisClient() redis.Provider {
	return &mockRedisClient{}
}

func StubVersionInfo() version.Info {
	return version.Info{}
}

func StubMetricProvider() *metrics.Provider {
	return &metrics.Provider{}
}

func StubLogger() log.Factory {
	return log.Factory{}
}

// func TestNew(t *testing.T) {
// 	type args struct {
// 		healthRepo     healthrepo.Store
// 		redisClient    redis.Provider
// 		versionInfo    version.Info
// 		metricProvider *metrics.Provider
// 		logger         log.Factory
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want Service
// 	}{
// 		{
// 			name: "WTF",
// 			args: args{
// 				healthRepo:     MockHealthRepo(),
// 				redisClient:    MockRedisClient(),
// 				versionInfo:    StubVersionInfo(),
// 				metricProvider: StubMetricProvider(),
// 				logger:         StubLogger(),
// 			},
// 			want: &service{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := New(tt.args.healthRepo, tt.args.redisClient, tt.args.versionInfo, tt.args.metricProvider, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_service_GetAPIUpTime(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx         context.Context
		versionInfo version.Info
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   healthmodels.APIUpTime
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			if got := svc.GetAPIUpTime(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetAPIUpTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetMemoryStats(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    healthmodels.MemoryStats
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			got, err := svc.GetMemoryStats(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetMemoryStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetMemoryStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetCacheHealth(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    healthmodels.CacheHealth
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			got, err := svc.GetCacheHealth(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetCacheHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetCacheHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetDatabaseHealth(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    healthmodels.DBHealth
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			got, err := svc.GetDatabaseHealth(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetDatabaseHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetDatabaseHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetAPIHealth(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    healthmodels.APIHealth
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			got, err := svc.GetAPIHealth(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetAPIHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetAPIHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetFullHealth(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx         context.Context
		versionInfo version.Info
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    healthmodels.TokenSvcHealth
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			got, err := svc.GetFullHealth(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetFullHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetFullHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_apiHealth(t *testing.T) {
	type fields struct {
		logger         log.Factory
		repo           healthrepo.Store
		redisClient    redis.Provider
		versionInfo    version.Info
		metricProvider *metrics.Provider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   healthmodels.APIHealth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				logger:         tt.fields.logger,
				repo:           tt.fields.repo,
				redisClient:    tt.fields.redisClient,
				versionInfo:    tt.fields.versionInfo,
				metricProvider: tt.fields.metricProvider,
			}
			if got := svc.apiHealth(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.apiHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memUsage(t *testing.T) {
	tests := []struct {
		name string
		want healthmodels.MemoryStats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := memStats(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}
