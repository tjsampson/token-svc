sum(irate(log_level_total[5m])) by (level)

sum(irate(http_request_total[5m])) by (method)

sum(irate(http_response_total[5m])) by (code)

sum(irate(http_request_total[5m])) by (uri)

sum(irate(container_memory_usage_bytes{name!=""}[5m])) by (name)

sum(irate(container_network_transmit_bytes_total{name!=""}[5m])) by (name)
