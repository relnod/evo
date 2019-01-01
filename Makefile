WD=$(shell pwd)

build:
	cd cmd/evod/ && go build -o ${WD}/out/evod
	cd cmd/evoproxy/ && go build -o ${WD}/out/evoproxy
	cd cmd/evoclient/ && go build -o ${WD}/out/evoclient

test:
	go test -v ./...

watch:
	modd -f scripts/modd.conf

plot:
	curl http://localhost:8080/stats > test.json
	python3 scripts/stats.py
