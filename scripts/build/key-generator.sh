#!/usr/bin/env bash

rm -rf /tmp/certs
mkdir /tmp/certs
openssl genrsa -out /tmp/certs/app.rsa 4096
openssl rsa -in /tmp/certs/app.rsa -pubout > /tmp/certs/app.rsa.pub