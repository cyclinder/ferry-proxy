package main

import (
	"context"
	"log"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/config"
	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/runtime"
)

func main() {
	ctx := context.Background()
	rt := runtime.NewRuntime(ctx)
	src := config.Source{
		EndpointCallback: rt.OnEndpoint,
		ClusterCallback:  rt.OnCluster,
		RouterCallback:   rt.OnRouter,
		ListenerCallback: rt.OnListener,
	}
	log.Println("start")
	err := src.Run("./config")
	if err != nil {
		panic(err)
	}
}
