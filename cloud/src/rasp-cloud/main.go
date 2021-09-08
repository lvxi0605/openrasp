//Copyright 2021-2021 corecna Inc.

package main

import (
	"rasp-cloud/controllers"
	_ "rasp-cloud/controllers"
	"rasp-cloud/environment"
	_ "rasp-cloud/environment"
	_ "rasp-cloud/filter"
	_ "rasp-cloud/models"
	"rasp-cloud/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.Listen.Graceful = true
	routers.InitRouter()
	beego.ErrorController(&controllers.ErrorController{})
	if environment.StartBeego {
		beego.Run()
	}
}
