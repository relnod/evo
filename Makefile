WD=$(shell pwd)

build:
	cd cmd/evod/ && go build -o ${WD}/out/evod
	cd cmd/evoproxy/ && go build -o ${WD}/out/evoproxy
	cd cmd/evoclient/ && go build -o ${WD}/out/evoclient

watch:
	modd -f scripts/modd.conf
