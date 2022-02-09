#! /bin/bash

# go build -o sensors cmd/web/*.go
# ./sensors

go build -o sensor cmd/web/*.go && ./sensor -production=false -cache=false -dbname=XXX -dbuser=XXX
