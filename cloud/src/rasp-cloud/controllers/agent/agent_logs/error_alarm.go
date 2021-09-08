//Copyright 2021-2021 corecna Inc.

package agent_logs

import (
	"rasp-cloud/conf"
	"rasp-cloud/controllers"
	"rasp-cloud/models/logs"
	"time"
)

type ErrorController struct {
	controllers.BaseController
}

// @router / [post]
func (o *ErrorController) Post() {
	var alarms []map[string]interface{}
	o.UnmarshalJson(&alarms)
	count := 0
	if conf.AppConfig.ErrorLogEnable {
		for _, alarm := range alarms {
			alarm["@timestamp"] = time.Now().UnixNano() / 1000000
			err := logs.AddErrorAlarm(alarm)
			if err == nil {
				count++
			}
		}
	}
	o.Serve(map[string]uint64{"count": uint64(count)})
}
