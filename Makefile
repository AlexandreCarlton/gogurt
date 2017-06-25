
BINARY := gogurt

all: build
.PHONY: all

deps:
	go get github.com/ulikunitz/xz

build: deps
	go build cmd/gogurt.go
.PHONY: build

run:
	go run cmd/gogurt.go
.PHONY: run

build-image:
	docker build --tag='centos-static' --build-arg user="${USER}" .
.PHONY: build-image

run-container:
	docker run --rm -it -v `pwd`:/home/$(USER) centos-static
.PHONY: build-image


