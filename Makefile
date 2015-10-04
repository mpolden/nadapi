NAME=nadapi

all: deps test lint install

deps:
	go get -d -v

fmt:
	go fmt ./...

lint:
	./lint.sh

test:
	go test ./...

install:
	go install
