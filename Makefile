.PHONY := all build install lint prepare test

GOPATH=$(PWD)

all: prepare lint build test install

build:
	(cd src/github.com/canterberry/goredis; GOPATH="$(GOPATH)" go build)

install:
	(cd src/github.com/canterberry/goredis; GOPATH="$(GOPATH)" go get)

lint:
	(cd src/github.com/canterberry/goredis; GOPATH="$(GOPATH)" go fmt)

prepare:
	(cd src/github.com/canterberry/goredis; GOPATH="$(GOPATH)" go get)

test:
	(cd src/github.com/canterberry/goredis; GOPATH="$(GOPATH)" go test)
