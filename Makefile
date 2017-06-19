MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
HOOKS:=pre-commit
LDFLAGSW:=-ldflags "-X github.com/autodidaddict/go-shopping/warehouse/internal/platform/config.Version=${VERSION}"
LDFLAGSC:=-ldflags "-X github.com/autodidaddict/go-shopping/catalog/internal/platform/config.Version=${VERSION}"
LDFLAGSS:=-ldflags "-X github.com/autodidaddict/go-shopping/shipping/internal/platform/config.Version=${VERSION}"

default: run

clean:
	rm -rf ./coverage.out ./coverage-all.out ./warehouse/cmd/warehoused/warehoused

cover: test
	go tool cover -html=coverage-all.out

build: clean
	@echo Building Warehouse Service...
	@cd warehouse/cmd/warehoused && CGO_ENABLED=0 go build ${LDFLAGSW} -a -installsuffix cgo -o warehoused main.go
	@echo Building Catalog Service...
	@cd catalog/cmd/catalogd && CGO_ENABLED=0 go build ${LDFLAGSC} -a -installsuffix cgo -o catalogd main.go
	@echo Building Shipping Service...
	@cd shipping/cmd/shippingd && CGO_ENABLED=0 go build ${LDFLAGSS} -a -installsuffix cgo -o shippingd main.go

catalog-proto:
	@cd catalog/proto && protoc --go_out=plugins=micro:. catalog.proto
warehouse-proto:
	@cd warehouse/proto && protoc --go_out=plugins=micro:. warehouse.proto
shipping-proto:
	@cd shipping/proto && protoc --go_out=plugins=micro:. shipping.proto

proto: shipping-proto catalog-proto warehouse-proto
	@echo All Protobufs Regenerated

hooks:
	cd .git/hooks
	$(foreach hook,$(HOOKS), ln -s -f ../../.git-hooks/${hook} .git/hooks/${hook};)

unhook:
	$(foreach hook,($HOOKS), unlink .git/hooks/${hook};)


