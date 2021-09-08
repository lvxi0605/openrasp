//Copyright 2021-2021 corecna Inc.

package agent_logs

import (
	"rasp-cloud/controllers"
	"rasp-cloud/models/logs"
	"time"
)

// Operations about policy alarm message
type PolicyAlarmController struct {
	controllers.BaseController
}

// @router / [post]
func (o *PolicyAlarmController) Post() {
	var alarms []map[string]interface{}
	o.UnmarshalJson(&alarms)
	count := 0
	for _, alarm := range alarms {
		alarm["@timestamp"] = time.Now().UnixNano() / 1000000
		err := logs.AddPolicyAlarm(alarm)
		if err == nil {
			count++
		}
	}
	o.Serve(map[string]uint64{"count": uint64(count)})
}
