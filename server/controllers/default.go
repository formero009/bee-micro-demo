package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	log.Println("get请求被调用")
	message := c.Ctx.Input.Param(":message")
	time.Sleep(200 * time.Millisecond)
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = Message{Method: c.Ctx.Request.Method, Message: message}
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
