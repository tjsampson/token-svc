#!/usr/bin/env bash

while sleep 1; do curl -d '{"email":"email@email.com","password":"somethingsupersecret", "confirm_password":"somethingsupersecret"}' -X POST https://dev.homerow.tech/api/v2/register; done
