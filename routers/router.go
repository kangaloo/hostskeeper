package routers

import (
	"github.com/astaxie/beego"
	"hostskeeper/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
