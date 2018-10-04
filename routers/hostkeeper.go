package routers

import (
	"github.com/astaxie/beego"
	"hostskeeper/controllers"
)

func init() {
	beego.Router("/hosts/add/id", &controllers.AddHosts{}, "*:AddByIDs")
	beego.Router("/hosts/add/ip", &controllers.AddHosts{}, "*:AddByIPs")
	beego.Router("/hosts/list/all", &controllers.Host{}, "*:ListAll")
	beego.Router("/hosts/list/init", &controllers.Host{}, "*:ListInit")
	beego.Router("/hosts/get/ip", &controllers.Host{}, "*:GetByIp")

//	setHostName setStatus getStatus ...
}
