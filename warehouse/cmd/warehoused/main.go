package main

import (
	"github.com/autodidaddict/go-shopping/warehouse/internal/platform/config"
	"github.com/autodidaddict/go-shopping/warehouse/internal/service"
	"github.com/autodidaddict/go-shopping/warehouse/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"time"
)

func main() {
	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	warehouse.RegisterWarehouseHandler(svc.Server(), service.NewWarehouseService())

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
