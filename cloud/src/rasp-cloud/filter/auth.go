//Copyright 2021-2021 corecna Inc.

package filter

import (
	"net/http"
	"rasp-cloud/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "X-CoreRASP-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.InsertFilter("/v1/agent/*", beego.BeforeRouter, authAgent)
	beego.InsertFilter("/v1/api/*", beego.BeforeRouter, authApi)
	beego.InsertFilter("/v1/iast/auth", beego.BeforeRouter, authAgent)
	beego.InsertFilter("/v1/iast/version", beego.BeforeRouter, authAgent)
	beego.InsertFilter("/v1/user/islogin", beego.BeforeRouter, authApi)
	beego.InsertFilter("/v1/user/default", beego.BeforeRouter, authApi)
}

func authAgent(ctx *context.Context) {
	appId := ctx.Input.Header("X-CoreRASP-AppID")
	appSecret := ctx.Input.Header("X-CoreRASP-AppSecret")
	app, err := models.GetAppById(appId)
	if appId == "" || err != nil || app == nil || appSecret != app.Secret {
		ctx.Output.JSON(map[string]interface{}{
			"status": http.StatusUnauthorized, "description": http.StatusText(http.StatusUnauthorized)},
			false, false)
	}
}

func authApi(ctx *context.Context) {
	cookie := ctx.GetCookie(models.AuthCookieName)
	if has, err := models.HasCookie(cookie); !has || err != nil {
		token := ctx.Input.Header(models.AuthTokenName)
		if has, err = models.HasToken(token); !has || err != nil {
			ctx.Output.JSON(map[string]interface{}{
				"status": http.StatusUnauthorized, "description": http.StatusText(http.StatusUnauthorized)},
				false, false)
			panic("")
		}
	}
}
