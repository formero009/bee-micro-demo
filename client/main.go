package main

import (
	config "client-demo/conf"
	_ "client-demo/routers"

	// wrappers "client-demo/wrappers"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	//read config
	info := config.GetInfo()

	//hystrix
	// hystrixClient := hystrix.NewClientWrapper(hystrix.WithFilter(func(c context.Context, e error) bool {
	// 	if e != nil {
	// 		fmt.Printf("e: %v\n", e)
	// 	}
	// 	return false
	// }))

	//conf beego
	beego.BConfig.CopyRequestBody = true
	//consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{info.Consul.Address}
	})
	srv := httpServer.NewServer(
		server.Name(info.Server.Name),
	)

	if err := srv.Handle(srv.NewHandler(beego.BeeApp.Handlers)); err != nil {
		logs.Error(err.Error())
		return
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Address(info.Server.Address),
		micro.Registry(reg),
	)

	service.Init()
	if err := service.Run(); err != nil {
		logs.Error("init service err")
		return
	}
}
