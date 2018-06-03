REPOSITORY?=relnod
NAME?=evo-base
TAG?=latest

.PHONY: all
all: build

.PHONY: build
build:
	docker build --tag=${REPOSITORY}/${NAME}:${TAG} .
