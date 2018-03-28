WORK_DIR=$(shell pwd)
OUT_DIR=${WORK_DIR}/out
STATIC_WEB_DIR=${WORK_DIR}/cmd/evo-web/static

${shell mkdir -p ${OUT_DIR}}

all: build

build: build-server build-web

build-server:
	cd cmd/evo-server && \
	go get ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${OUT_DIR}/evo-server

build-web: build-web-client build-web-app

build-web-client:
	cd cmd/evo-client && \
	go get -d -tags=js ./... && \
	gopherjs build -v -o ${STATIC_WEB_DIR}/evo-client.js

build-web-app:
	cd cmd/evo && \
	go get -d -tags=js ./... && \
	gopherjs build -v -o ${STATIC_WEB_DIR}/evo-app.js

clean:
	rm -rf ${OUT_DIR}

.PHONY: all build build-server build-web build-web-client build-web-app clean
