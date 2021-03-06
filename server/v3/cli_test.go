package v3

import (
	"context"
	httpClient "go-micro-demo/v3/httpclient"
	"testing"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
)

// var info config.Info

func TestHttpCli(t *testing.T) {
	//read config
	// info := config.GetInfo()

	//get service reg
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	//get service selector
	s := selector.NewSelector(selector.Registry(reg), selector.SetStrategy(selector.RoundRobin))

	//new http client
	c := httpClient.NewClient(client.Selector(s), client.DialTimeout(time.Second*10), client.RequestTimeout(time.Second*10))

	doGetRequest(t, c)
	doPostRequest(t, c)
}

func doGetRequest(t *testing.T, c client.Client) {
	request := c.NewRequest("go-micro-demo", "GET:/demo/hello/for-test/get", nil, client.WithContentType("application/json"))
	var response Resp
	if err := c.Call(context.Background(), request, &response); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("do get request success")
}

func doPostRequest(t *testing.T, c client.Client) {

	req := struct {
		Name string
		Age  int
	}{"jzd", 123}
	request := c.NewRequest("go-micro-demo", "POST:/demo/hello/for-test/post", req, client.WithContentType("application/json"))
	var response Resp
	if err := c.Call(context.Background(), request, &response); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("do post request success")
}

//msg
type Resp struct {
	Method  string
	Message string
}
