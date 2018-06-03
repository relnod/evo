REPOSITORY?=relnod
NAME?=evo-dev
TAG?=latest

.PHONY: all
all: build

.PHONY: build
build:
	docker build --tag=${REPOSITORY}/${NAME}:${TAG} .
