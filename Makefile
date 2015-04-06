NAME=nadapi
UID=$(shell id -u)

all: deps test lint install

deps:
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get -d -v

fmt:
	go fmt ./...

lint:
	go tool vet -copylocks=false $(PWD)
	golint ./...

test:
	go test ./...

install:
	go install

docker-build:
	docker run --rm -v $(PWD):/usr/src/$(NAME) -w /usr/src/$(NAME) \
		golang:latest /bin/sh -c \
		'go get -d -v && go build -v && chown $(UID):$(UID) $(NAME)'
