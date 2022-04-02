package main

import (
	"context"
	httpClient "go-micro-demo/v3/httpclient"
	"log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
)

// var info config.Info

func main() {
	//read config
	// info := config.GetInfo()

	//get service reg
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	//get service selector
	s := selector.NewSelector(
		selector.Registry(reg),
		selector.SetStrategy(selector.RoundRobin),
	)

	//配置服务的超时操作
	// configA := hystrix.CommandConfig{
	// 	MaxConcurrentRequests: 50,
	// 	Timeout:               2000000000,
	// }

	// hystrix.ConfigureCommand("go-micro-demo", configA)

	//将hystrix包装进客户端中
	// hystrix.ConfigureDefault(hystrix.CommandConfig{Timeout: 3000, MaxConcurrentRequests: 1})
	c := httpClient.NewClient(
		client.Selector(s),
		client.DialTimeout(time.Second*10),
		client.RequestTimeout(time.Second*10),
		// client.Wrap(hystrix.NewClientWrapper()),
	)

	doGetRequest(c)
	doPostRequest(c)
}

//hystrix 自定义
func doGetRequest(c client.Client) {
	request := c.NewRequest("go-micro-demo", "GET:/demo/hello/for-test/get", nil, client.WithContentType("application/json"))
	var response Resp
	err := hystrix.Do("test", func() error {
		err := c.Call(context.Background(), request, &response)
		return err
	}, nil)
	if err != nil {
		log.Fatalln(err.Error())
		return
	} else {
		log.Println("do get request success")
	}

}

func doPostRequest(c client.Client) {
	req := struct {
		Name string
		Age  int
	}{"jzd", 123}
	request := c.NewRequest("go-micro-demo", "POST:/demo/hello/for-test/post", req, client.WithContentType("application/json"))
	var response Resp
	if err := c.Call(context.Background(), request, &response); err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("do post request success")
}

//msg
type Resp struct {
	Method  string
	Message string
}
