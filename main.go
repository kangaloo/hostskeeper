package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "hostskeeper/routers"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	beego.Run()
}
