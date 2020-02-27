package healthmodels

import "github.com/tjsampson/token-svc/pkg/version"

// CacheHealth contains the redis cache health stats
type CacheHealth struct {
	Ping string `json:"ping"`
}

// APIUpTime represents the Server Uptime
type APIUpTime struct {
	Seconds float64 `json:"seconds"`
	Minutes float64 `json:"minutes"`
	Hours   float64 `json:"hours"`
	Days    float64 `json:"days"`
	Weeks   float64 `json:"weeks"`
}

// APIHealth contains status/health of the API
type APIHealth struct {
	Status  string       `json:"status"`
	Version version.Info `json:"version"`
	Uptime  APIUpTime    `json:"uptime"`
}

// DBHealth contains status/health of the DataBase
type DBHealth struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// TokenSvcHealth contains status/health of the TokenSvc Application Stack
type TokenSvcHealth struct {
	APIHealth   `json:"api_health"`
	DBHealth    `json:"db_health"`
	MemoryStats `json:"memory_stats"`
	CacheHealth `json:"cache_health"`
}

// MemoryStats holds the memory stats
type MemoryStats struct {
	Allocation      uint64 `json:"allocation"`
	TotalAllocation uint64 `json:"total_allocation"`
	Sys             uint64 `json:"sys"`
	NumGC           uint32 `json:"num_gc"`
	GoRoutines      int    `json:"go_routines"`
}
