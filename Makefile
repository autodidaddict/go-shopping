MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
HOOKS:=pre-commit
LDFLAGSW:=-ldflags "-X github.com/autodidaddict/go-shopping/warehouse/internal/platform/config.Version=${VERSION}"
LDFLAGSC:=-ldflags "-X github.com/autodidaddict/go-shopping/catalog/internal/platform/config.Version=${VERSION}"
LDFLAGSS:=-ldflags "-X github.com/autodidaddict/go-shopping/shipping/internal/platform/config.Version=${VERSION}"

default: run

test:
	@go test -v ./...

clean:
	@rm -rf ./coverage.out ./coverage-all.out ./warehouse/cmd/warehoused/warehoused

warehouse-lint:
	@golint -set_exit_status warehouse/internal/... warehouse/cmd/...

warehouse: clean warehouse-lint
	@echo Building Warehouse Service...
	@cd warehouse/cmd/warehoused && CGO_ENABLED=0 go build ${LDFLAGSW} -a -installsuffix cgo -o warehoused main.go

catalog-lint:
	@golint -set_exit_status catalog/internal/... catalog/cmd/...

catalog: clean catalog-lint
	@echo Building Catalog Service...
	@cd catalog/cmd/catalogd && CGO_ENABLED=0 go build ${LDFLAGSC} -a -installsuffix cgo -o catalogd main.go

shipping-lint:
	@golint -set_exit_status shipping/internal/... shipping/cmd/...

shipping: clean shipping-lint
	@echo Building Shipping Service...
	@cd shipping/cmd/shippingd && CGO_ENABLED=0 go build ${LDFLAGSS} -a -installsuffix cgo -o shippingd main.go

all: warehouse catalog shipping

catalog-proto:
	@cd catalog/proto && protoc --go_out=plugins=micro:. catalog.proto
warehouse-proto:
	@cd warehouse/proto && protoc --go_out=plugins=micro:. warehouse.proto
shipping-proto:
	@cd shipping/proto && protoc --go_out=plugins=micro:. shipping.proto

proto: shipping-proto catalog-proto warehouse-proto
	@echo All Protobufs Regenerated

hooks:
	chmod 755 .git-hooks/*
	cd .git/hooks
	$(foreach hook,$(HOOKS), ln -s -f ../../.git-hooks/${hook} .git/hooks/${hook};)

unhook:
	$(foreach hook,($HOOKS), unlink .git/hooks/${hook};)


