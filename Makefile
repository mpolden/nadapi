NAME=nadapi
UID=$(shell id -u)

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

docker-build:
	docker run --rm -v $(PWD):/usr/src/$(NAME) -w /usr/src/$(NAME) \
		golang:latest /bin/sh -c \
		'go get -d -v && go build -v && chown $(UID):$(UID) $(NAME)'
