package test

import (
	"client-demo/wrappers"
	"log"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/web"
	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
)

func main() {
	r := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})

	ginRouter := gin.Default()
	ginRouter.Handle("GET","/test",func(ctx *gin.Context) {
		ctx.String(200,"user api")
	})

	app := beego.NewApp

	// Create service
	webSvc := web.NewService(
		web.Name("testtest"),
		web.Address(":8081"),
		web.Registry(r),
		web.Handler(app.),
	)


	microSrv := micro.NewService(
		micro.Address(":9999"),
		micro.WrapClient(wrappers.NewlogWrap),
		micro.Registry(webSvc),
	)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
