//Copyright 2021-2021 corecna Inc.

package filter

import (
	"os"
	"rasp-cloud/conf"
	"rasp-cloud/tools"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

var (
	accessLogger *logs.BeeLogger
)

func init() {
	initAccessLogger()
	beego.InsertFilter("/*", beego.BeforeRouter, logAccess)
	beego.InsertFilter("/", beego.BeforeStatic, handleStatic)
}

func logAccess(ctx *context.Context) {
	var cont string
	cont += "[T]" + formatTime(time.Now().Unix(), "15:04:05") + " " + ctx.Input.Method() + " " +
		ctx.Input.Site() + ctx.Input.URI() + " - [I]" + ctx.Input.IP() + " | [U]" + ctx.Input.UserAgent()
	if ctx.Input.Referer() != "" {
		cont += "[F]" + ctx.Input.Referer()
	}
	if conf.AppConfig.RequestBodyEnable {
		body := ctx.Input.RequestBody
		cont += " - [B]" + string(body)
	}
	accessLogger.Info(cont)
}

func handleStatic(ctx *context.Context) {
	ctx.Output.Header("Cache-Control", "no-cache, no-store, max-age=0")
}

func formatTime(timestamp int64, format string) (times string) {
	tm := time.Unix(timestamp, 0)
	times = tm.Format(format)
	return
}

func initAccessLogger() {
	//var logPathSplit []string
	logFileName := "/access.log"
	maxSize := strconv.FormatInt(conf.AppConfig.LogMaxSize, 10)
	maxDays := strconv.Itoa(conf.AppConfig.LogMaxDays)
	logPath := conf.AppConfig.LogPath
	logAccessPath := conf.AppConfig.LogPath + "/access"
	// 判断后缀名称
	//if strings.HasSuffix(logPath, ".log") {
	//	logPathSplit = strings.Split(logPath, "/")
	//	logFileName = "/" + logPathSplit[len(logPathSplit) - 1]
	//	logPathSplitNoLogFileName := logPathSplit[:len(logPathSplit) - 1]
	//	logPath = strings.Join(logPathSplitNoLogFileName, "/")
	//}
	if isExists, _ := tools.PathExists(logPath); !isExists {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			tools.Panic(tools.ErrCodeLogInitFailed, "failed to create "+logPath+" dir", err)
		}
	}
	accessLogger = logs.NewLogger()
	accessLogger.EnableFuncCallDepth(true)
	accessLogger.SetLogFuncCallDepth(4)
	logAccessPath += logFileName
	err := accessLogger.SetLogger(logs.AdapterFile,
		`{"filename":"`+logAccessPath+`","daily":true,"maxdays":`+maxDays+`,"perm":"0777","maxsize": `+maxSize+`}`)
	if err != nil {
		tools.Panic(tools.ErrCodeLogInitFailed, "failed to init access log", err)
	}
}
