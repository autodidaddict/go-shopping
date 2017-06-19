package main

import (
	"github.com/autodidaddict/go-shopping/catalog/internal/platform/config"
	"github.com/autodidaddict/go-shopping/catalog/internal/service"
	"github.com/autodidaddict/go-shopping/catalog/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"time"
	"github.com/autodidaddict/go-shopping/catalog/internal/platform/redis"
)

func main() {
	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	redisCatalogRepository := redis.NewRedisRepository()
	catalog.RegisterCatalogHandler(svc.Server(), service.NewCatalogService(redisCatalogRepository))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
