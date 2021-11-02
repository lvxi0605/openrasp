//Copyright 2021-2021 corecna Inc.

package routers

import (
	"rasp-cloud/conf"
	"rasp-cloud/controllers"
	"rasp-cloud/controllers/agent"
	"rasp-cloud/controllers/agent/agent_logs"
	"rasp-cloud/controllers/api"
	"rasp-cloud/controllers/api/fore_logs"
	"rasp-cloud/controllers/iast"
	"rasp-cloud/tools"

	"github.com/astaxie/beego"
)

func InitRouter() {
	agentNS := beego.NewNamespace("/agent",
		beego.NSNamespace("/heartbeat",
			beego.NSInclude(
				&agent.HeartbeatController{},
			),
		),
		beego.NSNamespace("/log",
			beego.NSNamespace("/attack",
				beego.NSInclude(
					&agent_logs.AttackAlarmController{},
				),
			),
			beego.NSNamespace("/policy",
				beego.NSInclude(
					&agent_logs.PolicyAlarmController{},
				),
			),
			beego.NSNamespace("/error",
				beego.NSInclude(
					&agent_logs.ErrorController{},
				),
			),
		),
		beego.NSNamespace("/rasp",
			beego.NSInclude(
				&agent.RaspController{},
			),
		),
		beego.NSNamespace("/report",
			beego.NSInclude(
				&agent.ReportController{},
			),
		),
		beego.NSNamespace("/dependency",
			beego.NSInclude(
				&agent.DependencyController{},
			),
		),
		beego.NSNamespace("/crash",
			beego.NSInclude(
				&agent.CrashController{},
			),
		),
	)
	foregroudNS := beego.NewNamespace("/api",

		beego.NSNamespace("/plugin",
			beego.NSInclude(
				&api.PluginController{},
			),
		),
		beego.NSNamespace("/log",
			beego.NSNamespace("/attack",
				beego.NSInclude(
					&fore_logs.AttackAlarmController{},
				),
			),
			beego.NSNamespace("/policy",
				beego.NSInclude(
					&fore_logs.PolicyAlarmController{},
				),
			),
			beego.NSNamespace("/error",
				beego.NSInclude(
					&fore_logs.ErrorController{},
				),
			),
			beego.NSNamespace("/crash",
				beego.NSInclude(
					&fore_logs.CrashController{},
				),
			),
		),
		beego.NSNamespace("/app",
			beego.NSInclude(
				&api.AppController{},
			),
		),
		beego.NSNamespace("/display",
			beego.NSInclude(
				&api.DisplayController{},
			),
		),
		beego.NSNamespace("/rasp",
			beego.NSInclude(
				&api.RaspController{},
			),
		),
		beego.NSNamespace("/strategy",
			beego.NSInclude(
				&api.StrategyController{},
			),
		),
		beego.NSNamespace("/token",
			beego.NSInclude(
				&api.TokenController{},
			),
		),
		beego.NSNamespace("/report",
			beego.NSInclude(
				&api.ReportController{},
			),
		),
		beego.NSNamespace("/operation",
			beego.NSInclude(
				&api.OperationController{},
			),
		),
		beego.NSNamespace("/server",
			beego.NSInclude(
				&api.ServerController{},
			),
		),
		beego.NSNamespace("/dependency",
			beego.NSInclude(
				&api.DependencyController{},
			),
		),
	)
	iastNS := beego.NewNamespace("/iast",
		beego.NSInclude(
			&iast.WebsocketController{},
		),
		beego.NSInclude(
			&iast.IastController{},
		),
	)
	userNS := beego.NewNamespace("/user", beego.NSInclude(&api.UserController{}))
	pingNS := beego.NewNamespace("/ping", beego.NSInclude(&controllers.PingController{}))
	versionNS := beego.NewNamespace("/version", beego.NSInclude(&controllers.GeneralController{}))
	ns := beego.NewNamespace("/v1")
	ns.Namespace(pingNS)
	ns.Namespace(versionNS)
	startType := *conf.AppConfig.Flag.StartType
	if startType == conf.StartTypeForeground {
		ns.Namespace(foregroudNS, agentNS, userNS, iastNS)
	} else if startType == conf.StartTypeAgent {
		ns.Namespace(agentNS)
	} else if startType == conf.StartTypeDefault {
		ns.Namespace(foregroudNS, agentNS, userNS, iastNS)
	} else {
		tools.Panic(tools.ErrCodeStartTypeNotSupport, "Unknown -type parameter provided: "+startType, nil)
	}
	if startType == conf.StartTypeForeground || startType == conf.StartTypeDefault {
		beego.SetStaticPath("//", "dist")
	}
	beego.AddNamespace(ns)
}
