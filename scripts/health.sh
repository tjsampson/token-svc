#!/usr/bin/env bash

while sleep .5; do curl https://dev.homerow.tech/api/v2/health/ping; done
