# Use the official image to compile DreiTafel

.PHONY: build
build: deps
	go build -v

deps:
	go get

dbuild:
	docker run -it --rm -v "${GOPATH}":/go -v "${PWD}":/go/src/dreitafel -w /go/src/dreitafel golang go build

.PHONY: examples
examples:
	find examples/ -name "*.dot" -exec dot -Tpng -O {} \;
