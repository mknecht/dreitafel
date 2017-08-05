# Use the official image to compile DreiTafel

INDOCKER:=docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang

.PHONY: build
deps:
	$(INDOCKER) go get

build:
	$(INDOCKER) go build

dreitafel: *.go cmd/*.go
	$(INDOCKER) go build -o dreitafel cmd/main.go

.PHONY: examples
examples: dreitafel
	(cd examples && ./generate.sh)
