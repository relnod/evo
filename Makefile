WD=$(shell pwd)
# mkdir -p ${WD}/out/static

build:
	cd cmd/evod/ && go build -o ${WD}/out/evod
	cd cmd/evoproxy/ && go build -o ${WD}/out/evoproxy
	cd cmd/evoclient/ && go build -o ${WD}/out/evoclient
	cd cmd/evoclient/ && gopherjs build -o ${WD}/out/static/evoclient.js
	cp cmd/evoclient/index.html ${WD}/out/static/index.html

test:
	go test -v ./...

watch:
	modd -f scripts/modd.conf

plot:
	curl http://localhost:8080/stats > test.json
	python3 scripts/stats.py
