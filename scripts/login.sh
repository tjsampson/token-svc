#!/usr/bin/env bash

while sleep 1; do curl -d '{"username":"troy","password":"somethingsupersecret"}' -X POST https://dev.homerow.tech/api/v2/login; done
