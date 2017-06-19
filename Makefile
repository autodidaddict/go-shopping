MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
LDFLAGSW:=-ldflags "-X github.com/autodidaddict/go-shopping/warehouse/internal/platform/config.Version=${VERSION}"
LDFLAGSC:=-ldflags "-X github.com/autodidaddict/go-shopping/catalog/internal/platform/config.Version=${VERSION}"

default: run

clean:
	rm -rf ./coverage.out ./coverage-all.out ./warehouse/cmd/warehoused/warehoused

cover: test
	go tool cover -html=coverage-all.out

build: clean
	cd warehouse/cmd/warehoused && CGO_ENABLED=0 go build ${LDFLAGSW} -a -installsuffix cgo -o warehoused main.go
	cd catalog/cmd/catalogd && CGO_ENABLED=0 go build ${LDFLAGSC} -a -installsuffix cgo -o catalogd main.go

proto:
	cd warehouse/proto && protoc --go_out=plugins=micro:. warehouse.proto
	cd catalog/proto && protoc --go_out=plugins=micro:. catalog.proto