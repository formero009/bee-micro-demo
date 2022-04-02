package routers

import (
	"client-demo/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//route config
	ns := beego.NewNamespace("/demo",
		beego.NSNamespace("/hello",
			beego.NSInclude(
				&controllers.MainController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
