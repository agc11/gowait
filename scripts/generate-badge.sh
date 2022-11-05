#!/usr/bin/env bash

wget $(go test -race -coverprofile=coverage.out -covermode=atomic | grep -e "coverage" | awk '{print $2}' | sed 's/\./,/' | sed 's/%/%25/'  | awk '{print "https://badgen.net/badge/coverage/"$1"/green?icon=github"}') -O ./assets/coverage/coverage.svg