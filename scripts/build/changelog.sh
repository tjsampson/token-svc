#!/usr/bin/env bash

startTag=$1
endTag=$2

git log ${startTag}...${endTag} --pretty=format:'- [%ad] %an: [%s](http://github.com/tjsampson/token-svc/commit/%H)' >> CHANGELOG.md