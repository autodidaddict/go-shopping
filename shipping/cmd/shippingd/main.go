package main

import (
	"github.com/autodidaddict/go-shopping/shipping/internal/platform/config"
	"github.com/autodidaddict/go-shopping/shipping/internal/service"
	"github.com/autodidaddict/go-shopping/shipping/proto"
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

	shipping.RegisterShippingHandler(svc.Server(), service.NewShippingService())

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
