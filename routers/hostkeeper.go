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
	beego.Router("/hosts/add", &controllers.Host{}, "*:Add")
	//	curl "127.0.0.1:8080/hosts/add?ip=192.168.68.82&hostname=ecs7&cpu=2&mem=4&disk=200"

	//	setHostName setStatus getStatus ...
}
