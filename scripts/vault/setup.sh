#!/usr/bin/env bash

## CONFIG LOCAL ENV
echo "[*] Config local environment..."
export VAULT_ADDR=http://0.0.0.0:8200

mkdir data

## INIT VAULT
echo "[*] Init vault..."
vault operator init -address=${VAULT_ADDR} > ./data/keys.txt
export VAULT_TOKEN=$(grep 'Initial Root Token:' ./data/keys.txt | awk '{print substr($NF, 1, length($NF))}')

## UNSEAL VAULT
echo "[*] Unseal vault..."
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 1:' ./data/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 2:' ./data/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 3:' ./data/keys.txt | awk '{print $NF}')

## AUTH
echo "[*] Auth..."
vault login -address=${VAULT_ADDR} ${VAULT_TOKEN}

## CREATE USER
echo "[*] Create user... Remember to change the defaults!!"
vault auth enable userpass
vault policy write admin /config/admin.hcl
vault write -address=${VAULT_ADDR} auth/userpass/users/webui password=webui policies=admin

## CREATE BACKUP TOKEN
echo "[*] Create backup token..."
vault token create -address=${VAULT_ADDR} -display-name="backup_token" | awk '/token/{i++}i==2' | awk '{print "backup_token: " $2}' >>./data/keys.txt

# ENABLE KV /SECRET
echo "[*] Creating new mount point..."
vault secrets enable -path=secret/ kv

## READ/WRITE
# $ vault write -address=${VAULT_ADDR} secret/api-key value=12345678
# $ vault read -address=${VAULT_ADDR} secret/api-key
