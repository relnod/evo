WORK_DIR=$(shell pwd)
OUT_DIR=${WORK_DIR}/out
STATIC_WEB_DIR=${OUT_DIR}/static
REPO_DIR=github.com/relnod/evo

${shell mkdir -p ${OUT_DIR}}
${shell mkdir -p ${STATIC_WEB_DIR}}

.PHONY: all
all: build

.PHONY: dev
dev:
	docker-compose -f docker-compose.yml -f docker-compose-dev.yml up --build

.PHONY: run-desktio-client
run-desktop-client: build-desktop-client
	${OUT_DIR}/evo-client

.PHONY: build
build: build-server build-web

.PHONY: build-server
build-server:
	cd cmd/evo-server && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUT_DIR}/evo-server
	# cp cmd/evo-web/static/* ${STATIC_WEB_DIR}/

.PHONY: build-web
build-web: build-web-client build-web-app

.PHONY: build-web-client
build-web-client:
	cd cmd/evo-client && \
	gopherjs build -o ${STATIC_WEB_DIR}/evo-client.js

.PHONY: build-web-app
build-web-app:
	cd cmd/evo && \
	gopherjs build -o ${STATIC_WEB_DIR}/evo-app.js

.PHONY: build-desktop-client
build-desktop-client:
	go build -o ${OUT_DIR}/evo-client ${REPO_DIR}/cmd/evo-client

.PHONY: build-desktop-app
build-desktop-app:
	go build -o ${OUT_DIR}/evo-app ${REPO_DIR}/cmd/evo

.PHONY: clean
clean:
	rm -rf ${OUT_DIR}
