#!/bin/bash
docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang go run cmd/main.go "$*"
