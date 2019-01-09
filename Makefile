WD=$(shell pwd)

export GO111MODULE=on

build: clean
	cd cmd/evod/ && go build -o ${WD}/out/evod
	cd cmd/evoproxy/ && go build -o ${WD}/out/evoproxy
	cd cmd/evoclient/ && go build -o ${WD}/out/evoclient
	cd cmd/evoclient/ && gopherjs build -o ${WD}/out/static/evoclient.js
	cp cmd/evoclient/index.html ${WD}/out/static/index.html

test:
	# TODO: find a way to run tests in a package that uses opengl
	go test -v `go list ./... | grep -v cmd | grep -v graphics`

bench:
	go test -run=NONE -bench=. ./...

watch:
	modd -f scripts/modd.conf

plot:
	curl http://localhost:8080/stats > test.json
	python3 scripts/stats.py

clean:
	rm -rf out
