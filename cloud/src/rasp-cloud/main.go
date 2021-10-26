//Copyright 2021-2021 corecna Inc.

package main

import (
	"log"
	"rasp-cloud/config"
	"rasp-cloud/controllers"
	"rasp-cloud/environment"
	_ "rasp-cloud/filter"
	_ "rasp-cloud/models"
	"rasp-cloud/routers"
	"rasp-cloud/tools"

	"github.com/astaxie/beego"
)

func init() {

}

func main() {

	checkflag := tools.CheckAppkeyAndSecret(config.TOMLConfig.AppKey, config.TOMLConfig.AppSecret)
	if !checkflag {
		log.Printf("appkey or appsecret is invalid")
		return
	}

	beego.BConfig.Listen.Graceful = true
	routers.InitRouter()
	beego.ErrorController(&controllers.ErrorController{})
	if environment.StartBeego {
		beego.Run()
	}
}
