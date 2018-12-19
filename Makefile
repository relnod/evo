WD=$(shell pwd)

build:
	cd cmd/evod/ && go build -o ${WD}/out/evod
	cd cmd/evoclient/ && gopherjs build main.go
	cd cmd/evoclient/ && go build -o ${WD}/out/evoclient
	cd cmd/evo/ && go build -o ${WD}/out/evo

watch:
	modd -f scripts/modd.conf
