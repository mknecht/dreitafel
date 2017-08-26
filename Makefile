# Use the official image to compile DreiTafel

GO_UBUNTU:=docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang
GO_ALPINE:=docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang:alpine

.PHONY: all
all: dreitafel dreitafel-web

deps:
	$(GO_UBUNTU) go get

dreitafel: *.go cmd/*.go
	$(GO_UBUNTU) go build -o dreitafel cmd/main.go

dreitafel-web: *.go web/*.go cmd/*.go
	$(GO_UBUNTU) go build -o dreitafel-web cmd/web.go

.PHONY: examples
examples: dreitafel
	(cd examples && ./generate.sh)

clean:
	rm -f dreitafel dreitafel-web

docker-image-latest:
	$(GO_ALPINE) go build -o docker/dreitafel cmd/main.go
	$(GO_ALPINE) go build -o docker/dreitafel-web cmd/web.go
	docker build -t muratk/dreitafel:latest docker/
