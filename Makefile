
BINARY := gogurt

all: build
.PHONY: all

deps:
	go get github.com/ulikunitz/xz

build: deps
	go build github.com/alexandrecarlton/gogurt/cmd/gogurt
.PHONY: build

build-image:
	docker build --tag='centos-static' --build-arg user="${USER}" .
.PHONY: build-image

run-container:
	docker run --rm -it -v `pwd`:/home/$(USER) centos-static
.PHONY: build-image


