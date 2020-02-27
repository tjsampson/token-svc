package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type api struct {
	ServiceName         string   `toml:"servicename"`
	ShutdownTimeoutSecs uint16   `toml:"shutdowntimeoutsecs"`
	ReadTimeOutSecs     uint16   `toml:"readtimeoutsecs"`
	WriteTimeOutSecs    uint16   `toml:"writetimeoutsecs"`
	IdleTimeOutSecs     uint16   `toml:"idletimeoutsecs"`
	TimeoutSecs         uint16   `toml:"timeoutsecs"`
	Port                string   `toml:"port"`
	MetricsPort         string   `toml:"metricsport"`
	AllowedMethods      []string `toml:"allowedmethods"`
	AllowedOrigins      []string `toml:"allowedorigins"`
	AllowedHeaders      []string `toml:"allowedheaders"`
	OpenEndPoints       []string `toml:"openendpoints"`
}

type token struct {
	AccessTokenLifeSpanMins             uint16 `toml:"accesstokenlifespanmins"`
	RefreshTokenLifeSpanMins            uint16 `toml:"refreshtokelifespanmins"`
	FailedLoginAttemptCacheLifeSpanMins uint16 `toml:"failedloginattemptcachelifespanmins"`
	FailedLoginAttemptsMax              uint16 `toml:"failedloginattemptsmax"`
	AuthPrivateKeyPath                  string `toml:"authprivatekeypath"`
	AuthPublicKeyPath                   string `toml:"authpublickeypath"`
	Issuer                              string `toml:"issuer"`
	AccessCacheKeyID                    string `tom:"accesscachekeyid"`
	FailedLoginCacheKeyID               string `toml:"failedlogincachekeyid"`
	RefreshCacheKeyID                   string `tom:"refreshcachekeyid"`
}

type db struct {
	User    string `toml:"user"`
	Pass    string `toml:"pass"`
	Host    string `toml:"host"`
	Port    string `toml:"port"`
	Name    string `toml:"name"`
	Timeout string `toml:"timeout"`
}

type cookie struct {
	LifeSpanDays    uint16 `toml:"lifespandays"`
	HashKey         string `toml:"hashkey"`
	BlockKey        string `toml:"blockkey"`
	Name            string `toml:"name"`
	Domain          string `toml:"domain"`
	KeyUserID       string `toml:"keyuserid"`
	KeyEmail        string `toml:"keyemail"`
	KeyJWTAccessID  string `toml:"keyjwtaccessid"`
	KeyJWTRefreshID string `toml:"keyjwtrefreshid"`
}

type cache struct {
	Host                          string `toml:"host"`
	Port                          string `toml:"port"`
	UserAccountLockedKeyID        string `toml:"useraccountlockedkeyid"`
	UserAccountLockedLifeSpanMins uint16 `toml:"useraccountlockedlifespanmins"`
}

type logger struct {
	Level            string   `toml:"level"`
	Encoding         string   `toml:"encoding"`
	OutputPaths      []string `toml:"outputpaths"`
	ErrorOutputPaths []string `toml:"erroroutputpaths"`
}

// Config the configuration struct for the service
type Config struct {
	API    api    `toml:"api"`
	Logger logger `toml:"logger"`
	Token  token  `toml:"token"`
	DB     db     `toml:"db"`
	Cookie cookie `toml:"cookie"`
	Cache  cache  `toml:"cache"`
}

// defConfig which is sane defaults for development purposes (local).
// Other environments should use the TOKEN_SVC_CONF environment variable
// ex:
// 		$ TOKEN_SVC_CONF=/path/to/token-svc/config.toml ./bin/token-svc
func defConfig() Config {
	return Config{
		API: api{
			ServiceName:         "token-svc",
			MetricsPort:         "4001",
			Port:                "4000",
			ShutdownTimeoutSecs: 120,
			IdleTimeOutSecs:     90,
			WriteTimeOutSecs:    30,
			ReadTimeOutSecs:     5,
			TimeoutSecs:         30,
			AllowedHeaders:      []string{"X-Requested-With", "X-Request-ID", "jaeger-debug-id", "Content-Type", "Authorization"},
			AllowedOrigins:      []string{"*"},
			AllowedMethods:      []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
			OpenEndPoints:       []string{"/login", "/health/ping", "/register"},
		},
		Logger: logger{
			Level:            "debug",
			Encoding:         "json",
			OutputPaths:      []string{"stdout", "/tmp/logs/tokensvc.logs"},
			ErrorOutputPaths: []string{"stderr"},
		},
		DB: db{
			User:    "postgres",
			Pass:    "postgres",
			Host:    "postgres",
			Port:    "5432",
			Name:    "postgres",
			Timeout: "30",
		},
		Token: token{
			AccessTokenLifeSpanMins:             30,    // half hour
			RefreshTokenLifeSpanMins:            10080, // 1 week
			FailedLoginAttemptCacheLifeSpanMins: 30,
			FailedLoginAttemptsMax:              5,
			AuthPrivateKeyPath:                  "/tmp/certs/app.rsa", // TODO: Let's read these in from Vault
			AuthPublicKeyPath:                   "/tmp/certs/app.rsa.pub",
			Issuer:                              "homerow.tech",
			AccessCacheKeyID:                    "token-access-user",
			RefreshCacheKeyID:                   "token-refresh-user",
			FailedLoginCacheKeyID:               "failed-login-user",
		},
		Cookie: cookie{
			LifeSpanDays:    7,
			HashKey:         "something-that-is-32-byte-secret",
			BlockKey:        "something-else-16-24-or-32secret",
			Name:            "homerow.tech",
			Domain:          "dev.homerow.tech",
			KeyUserID:       "id",
			KeyEmail:        "email",
			KeyJWTAccessID:  "jti-access",
			KeyJWTRefreshID: "jti-refresh",
		},
		Cache: cache{
			Host:                          "redis",
			Port:                          "6379",
			UserAccountLockedLifeSpanMins: 60,
			UserAccountLockedKeyID:        "account-locked-user",
		},
	}
}

// Read reads the config file and return the populated Config Struct
func Read() *Config {
	var cfg Config
	var cfgPath string

	cfgPath = os.Getenv("TOKEN_SVC_CONF")

	if len(cfgPath) <= 0 {
		log.Println("TOKEN_SVC_CONF env variable not set")
		log.Println("using default config value...")
		cfg = defConfig()
	} else {
		_, err := toml.DecodeFile(cfgPath, &cfg)

		if err != nil {
			log.Fatalf("unable to parse config: %v", err.Error())
		}
	}
	return &cfg
}
