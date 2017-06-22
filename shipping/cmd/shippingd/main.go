package main

import (
	"github.com/autodidaddict/go-shopping/shipping/internal/platform/config"
	"github.com/autodidaddict/go-shopping/shipping/internal/platform/redis"
	"github.com/autodidaddict/go-shopping/shipping/internal/service"
	"github.com/autodidaddict/go-shopping/shipping/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"time"
)

func main() {
	repo := redis.NewRedisRepository(":6379")
	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	shipping.RegisterShippingHandler(svc.Server(), service.NewShippingService(repo))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
