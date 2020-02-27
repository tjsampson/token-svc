#!/usr/bin/env bash


vault kv put secret/services/token-svc/config/db user=postgres pass=postgres

vault kv put secret/services/token-svc/config/cookie hashkey=c88324985aad39b9174487786eb73783 blockkey=94bcab6a2660bf33a8429501edc6dd0a

