[api]
servicename = "{{ key "services/token-svc/config/api/servicename" }}"
port = "{{ key "services/token-svc/config/api/port" }}"
metricsport = "{{ key "services/token-svc/config/api/metricsport" }}"
allowedmethods = {{ key "services/token-svc/config/api/allowedmethods" }}
allowedorigins = {{ key "services/token-svc/config/api/allowedorigins" }}
allowedheaders = {{ key "services/token-svc/config/api/allowedheaders" }}
openendpoints = {{ key "services/token-svc/config/api/openendpoints" }}
shutdowntimeoutsecs = {{ key "services/token-svc/config/api/shutdowntimeoutsecs" }}                 
idletimeoutsecs = {{ key "services/token-svc/config/api/idletimeoutsecs" }}                 
writetimeoutsecs = {{ key "services/token-svc/config/api/writetimeoutsecs" }}                   
readtimeoutsecs = {{ key "services/token-svc/config/api/readtimeoutsecs" }}                 
timeoutsecs = {{ key "services/token-svc/config/api/timeoutsecs" }}                 

[logger]
level = "{{ key "services/token-svc/config/logger/level" }}"
encoding = "{{ key "services/token-svc/config/logger/encoding" }}"
outputpaths = {{ key "services/token-svc/config/logger/outputpaths" }}
erroroutputpaths = {{ key "services/token-svc/config/logger/erroroutputpaths" }}

[db]
user = "{{ with secret "secret/services/token-svc/config/db" }}{{ .Data.user }}{{ end }}"
pass = "{{ with secret "secret/services/token-svc/config/db" }}{{ .Data.pass }}{{ end }}"
host = "{{ key "services/token-svc/config/db/host" }}"
port = "{{ key "services/token-svc/config/db/port" }}"
name = "{{ key "services/token-svc/config/db/name" }}"
timeout = "{{ key "services/token-svc/config/db/timeout" }}"

[token]
authprivatekeypath = "{{ key "services/token-svc/config/token/authprivatekeypath" }}"
authpublickeypath = "{{ key "services/token-svc/config/token/authpublickeypath" }}"
issuer = "{{ key "services/token-svc/config/token/issuer" }}"
accesstokenlifespanmins = {{ key "services/token-svc/config/token/accesstokenlifespanmins" }}
refreshtokenlifespanmins = {{ key "services/token-svc/config/token/refreshtokenlifespanmins" }}
accesscachekeyid = "{{ key "services/token-svc/config/token/accesscachekeyid" }}"
refreshcachekeyid = "{{ key "services/token-svc/config/token/refreshcachekeyid" }}"
failedlogincachekeyid = "{{ key "services/token-svc/config/token/failedlogincachekeyid" }}"
failedloginattemptcachelifespanmins = "{{ key "services/token-svc/config/token/failedloginattemptcachelifespanmins" }}"
failedloginattemptsmax = "{{ key "services/token-svc/config/token/failedloginattemptsmax" }}"

[cookie]
hashkey = "{{ with secret "secret/services/token-svc/config/cookie" }}{{ .Data.hashkey }}{{ end }}"
blockkey = "{{ with secret "secret/services/token-svc/config/cookie" }}{{ .Data.blockkey }}{{ end }}"
domain = "{{ key "services/token-svc/config/cookie/domain" }}"
name = "{{ key "services/token-svc/config/cookie/name" }}"
lifespandays = {{ key "services/token-svc/config/cookie/lifespandays" }}
keyuserid = "{{ key "services/token-svc/config/cookie/keyuserid" }}"
keyemail = "{{ key "services/token-svc/config/cookie/keyemail" }}"
keyjwtaccessid = "{{ key "services/token-svc/config/cookie/keyjwtaccessid" }}"
keyjwtrefreshid = "{{ key "services/token-svc/config/cookie/keyjwtrefreshid" }}"

[cache]
host = "{{ key "services/token-svc/config/cache/host" }}"
port = "{{ key "services/token-svc/config/cache/port" }}"
useraccountlockedkeyid = "{{ key "services/token-svc/config/cache/useraccountlockedkeyid" }}"
useraccountlockedlifespanmins = "{{ key "services/token-svc/config/cache/useraccountlockedlifespanmins" }}"
