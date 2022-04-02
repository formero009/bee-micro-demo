package main

import (
	"context"
	"fmt"
	config "go-micro-demo/conf/yaml_config"
	_ "go-micro-demo/routers"
	"log"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// //go-micro Wrapper中间件
// type ProductWrapper struct {
// 	client.Client
// }

// //Wrapper中间件的执行方法
// func (this *ProductWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
// 	cmdName := req.Service() + "." + req.Endpoint()
// 	config := hystrix.CommandConfig{
// 		Timeout: 1000,
// 	}
// 	hystrix.ConfigureCommand(cmdName, config)
// 	return hystrix.Do(cmdName, func() error {
// 		return this.Client.Call(ctx, req, rsp)
// 	}, func(e error) error {
// 		// DefaultProducts(rsp)
// 		return nil
// 	})
// }

// //初始化Wrapper
// func NewProductWrapper(c client.Client) client.Client {
// 	return &ProductWrapper{c}
// }

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	fmt.Printf("[Log Wrapper] ctx: %v service: %s method: %s\n", md, req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

func NewLogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[Log Wrapper] Before serving request method: %v", req.Endpoint())
		err := fn(ctx, req, rsp)
		log.Printf("[Log Wrapper] After serving request")
		return err
	}
}

func main() {
	//read config
	info := config.GetInfo()

	//conf beego
	beego.BConfig.CopyRequestBody = true
	//consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{info.Consul.Address}
	})
	srv := httpServer.NewServer(
		server.Name(info.Server.Name),
		server.Address("127.0.0.1:8010"),
	)

	if err := srv.Handle(srv.NewHandler(beego.BeeApp.Handlers)); err != nil {
		logs.Error(err.Error())
		return
	}

	service := micro.NewService(
		micro.Server(srv),
		micro.Address(info.Server.Address),
		micro.Registry(reg),
		//熔断
		micro.WrapHandler(NewLogWrapper),
		// micro.WrapClient(hystrix.NewClientWrapper()),
		//负载均衡
		// micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()

	if err := service.Run(); err != nil {
		logs.Error("init service err", err)
		return
	}
}
