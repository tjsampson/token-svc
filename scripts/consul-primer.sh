#!/usr/bin/env bash

consul kv put services/token-svc/config/api/servicename token-svc
consul kv put services/token-svc/config/api/metricsport 4001
consul kv put services/token-svc/config/api/port 4000
consul kv put services/token-svc/config/api/allowedmethods '["GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"]'
consul kv put services/token-svc/config/api/allowedorigins '["*"]'
consul kv put services/token-svc/config/api/allowedheaders '["X-Requested-With","X-Request-ID", "jaeger-debug-id", "Content-Type", "Authorization"]'
consul kv put services/token-svc/config/api/openendpoints '["/login", "/health/ping", "/register"]'
consul kv put services/token-svc/config/api/shutdowntimeoutsecs 120
consul kv put services/token-svc/config/api/idletimeoutsecs 90
consul kv put services/token-svc/config/api/writetimeoutsecs 30
consul kv put services/token-svc/config/api/readtimeoutsecs 5
consul kv put services/token-svc/config/api/timeoutsecs 30
consul kv put services/token-svc/config/logger/level debug
consul kv put services/token-svc/config/logger/encoding json
consul kv put services/token-svc/config/logger/outputpaths '["stdout", "/tmp/tokensvc.logs"]'
consul kv put services/token-svc/config/logger/erroroutputpaths '["stderr"]'
consul kv put services/token-svc/config/db/host postgres
consul kv put services/token-svc/config/db/port 5432
consul kv put services/token-svc/config/db/name postgres
consul kv put services/token-svc/config/db/timeout 5
consul kv put services/token-svc/config/token/authprivatekeypath '/tmp/certs/app.rsa'
consul kv put services/token-svc/config/token/authpublickeypath '/tmp/certs/app.rsa.pub'
consul kv put services/token-svc/config/token/issuer 'homerow.tech'
consul kv put services/token-svc/config/token/accesstokenlifespanmins 30
consul kv put services/token-svc/config/token/refreshtokenlifespanmins 10080
consul kv put services/token-svc/config/token/accesscachekeyid 'token-access-user'
consul kv put services/token-svc/config/token/refreshcachekeyid 'token-refresh-user'
consul kv put services/token-svc/config/token/failedlogincachekeyid 'failed-login-user'
consul kv put services/token-svc/config/token/failedloginattemptcachelifespanmins 30
consul kv put services/token-svc/config/cookie/domain 'dev.homerow.tech'
consul kv put services/token-svc/config/cookie/name 'homerow-auth'
consul kv put services/token-svc/config/cookie/lifespandays 7
consul kv put services/token-svc/config/cookie/keyuserid 'id'
consul kv put services/token-svc/config/cookie/keyemail 'email'
consul kv put services/token-svc/config/cookie/keyjwtaccessid 'jti-access'
consul kv put services/token-svc/config/cookie/keyjwtrefreshid 'jti-refresh'
consul kv put services/token-svc/config/cache/host 'redis'
consul kv put services/token-svc/config/cache/port '6379'
consul kv put services/token-svc/config/cache/useraccountlockedkeyid = "account-locked-user"
consul kv put services/token-svc/config/cache/useraccountlockedlifespanmins = 60