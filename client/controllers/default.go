package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	httpClient "client-demo/client"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	beego.Controller
}

type Message struct {
	Method  string
	Message string
}

type PostInfo struct {
	Name string
	Age  int
}

// @Title get test
// @Description get test
// @Success 200 success message
// @router /:message/get [get]
func (c *MainController) Get() {
	//模拟请求服务端 get请求
	// var resp httpClient.Resp
	log.Println("进入client 请求，准备向服务端发送get请求...")
	//获取链接中的message
	message := c.Ctx.Input.Param(":message")
	t, err := httpClient.DoGetRequest("go-micro-demo", "GET:/demo/hello/"+message+"/get", nil)
	if err != nil {
		log.Println("远程调用出错!")
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = Message{Method: c.Ctx.Request.Method, Message: err.Error()}
	} else {
		log.Println("远程调用成功")
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = Message{Method: c.Ctx.Request.Method, Message: t.Message}
	}
	c.ServeJSON()
}

// @Title post test
// @Description post test
// @Success 200 success message
// @router /:message/post [post]
func (c *MainController) Post() {
	var v PostInfo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error("get body error. %v", err)
		return
	}
	message := c.Ctx.Input.Param(":message")
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = Message{Method: c.Ctx.Request.Method, Message: message}
	c.ServeJSON()
}
