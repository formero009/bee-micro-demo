package v3

import (
	config "client-demo/conf"
	httpClient "client-demo/v3/httpclient"
	wrappers "client-demo/wrappers"
	"context"
	"log"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
)

func NewHttpClient() client.Client {

	info := config.GetInfo()

	//get service reg
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{info.Consul.Address}
	})

	//get service selector
	s := selector.NewSelector(selector.Registry(reg), selector.SetStrategy(selector.RoundRobin))

	//new http client
	client := httpClient.NewClient(client.Selector(s),
		client.DialTimeout(time.Second*10),
		client.RequestTimeout(time.Second*10),
		client.Wrap(wrappers.NewClientWrapper()),
		client.Wrap(wrappers.NewlogWrap),
	)

	return client
}

//"http-demo", "GET:/demo/hello/for-test/get", nil, client.WithContentType("application/json")
func DoGetRequest(service string, endpoint string, req interface{}, reqOpts ...client.RequestOption) (response Resp, err error) {
	c := NewHttpClient()

	//默认json
	reqOpts = append(reqOpts, client.WithContentType("application/json"))
	request := c.NewRequest(service, endpoint, req, reqOpts...)

	if err = c.Call(context.Background(), request, &response); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("DoGetRequest执行成功！")
	return response, nil
}

//"http-demo", "POST:/demo/hello/for-test/post", req, client.WithContentType("application/json")
func DoPostRequest(service string, endpoint string, req interface{}, reqOpts ...client.RequestOption) (response Resp, err error) {
	c := NewHttpClient()
	// req = struct {
	// 	Name string
	// 	Age  int
	// }{"jzd", 123}
	//默认json
	reqOpts = append(reqOpts, client.WithContentType("application/json"))
	request := c.NewRequest(service, endpoint, req, reqOpts...)

	if err = c.Call(context.Background(), request, &response); err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("do post request success")
	return response, err
}

//msg
type Resp struct {
	Method  string
	Message string
}
