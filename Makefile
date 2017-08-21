# Use the official image to compile DreiTafel

INDOCKER:=docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang

.PHONY: all
all: dreitafel dreitafel-web

deps:
	$(INDOCKER) go get

dreitafel: *.go cmd/*.go
	$(INDOCKER) go build -o dreitafel cmd/main.go

dreitafel-web: *.go web/*.go cmd/*.go
	$(INDOCKER) go build -o dreitafel-web cmd/web.go

.PHONY: examples
examples: dreitafel
	(cd examples && ./generate.sh)

clean:
	rm -f dreitafel dreitafel-web

docker-image-latest:
	cp dreitafel docker/
	cp dreitafel-web docker/
	docker build -t muratk/dreitafel:latest docker/
