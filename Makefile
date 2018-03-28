WORK_DIR=$(shell pwd)
OUT_DIR=${WORK_DIR}/out

${shell mkdir -p ${OUT_DIR}}

all: build

build: build-server build-web

build-server:
	cd cmd/evo-server && \
	go get ./... && \
	go build -o ${OUT_DIR}/evo-server

build-web: build-web-client build-web-app

build-web-client:
	cd cmd/evo-client && \
	go get -d -tags=js ./... && \
	gopherjs build -v -o ${OUT_DIR}/evo-client.js

build-web-app:
	cd cmd/evo && \
	go get -d -tags=js ./... && \
	gopherjs build -v -o ${OUT_DIR}/evo-app.js

clean:
	rm -rf ${OUT_DIR}

.PHONY: all build build-server build-web build-web-client build-web-app clean
