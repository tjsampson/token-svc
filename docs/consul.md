# Consul

Vault utilizes Consul for it's backend storage, and Consul is also used as a Service Registry. Consul-Template is also used to read configuration data from Consul and Vault

## Consul-Template

`consul-template -log-level trace -template "config.tpl:config.toml" -once`
