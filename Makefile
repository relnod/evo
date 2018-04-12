WORK_DIR=$(shell pwd)
OUT_DIR=${WORK_DIR}/out
STATIC_WEB_DIR=${OUT_DIR}/static

${shell mkdir -p ${OUT_DIR}}
${shell mkdir -p ${STATIC_WEB_DIR}}

all: build

dev:
	docker-compose -f docker-compose.yml -f docker-compose-dev.yml up --build

build: build-server build-web

build-server:
	cd cmd/evo-server && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUT_DIR}/evo-server
	# cp cmd/evo-web/static/* ${STATIC_WEB_DIR}/

build-web: dep-web build-web-client build-web-app

build-web-client:
	cd cmd/evo-client && \
	gopherjs build -o ${STATIC_WEB_DIR}/evo-client.js

build-web-app:
	cd cmd/evo && \
	gopherjs build -o ${STATIC_WEB_DIR}/evo-app.js

dep-web:
	go get -tags=js github.com/goxjs/glfw
	go get github.com/gopherjs/gopherjs
	go get github.com/gopherjs/gopherjs/js

clean:
	rm -rf ${OUT_DIR}

.PHONY: all build build-server build-web build-web-client build-web-app clean
