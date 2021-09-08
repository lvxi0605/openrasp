//Copyright 2021-2021 corecna Inc.

package agent_logs

import (
	"rasp-cloud/controllers"
	"rasp-cloud/models/logs"
	"time"
)

// Operations about attack alarm message
type AttackAlarmController struct {
	controllers.BaseController
}

// @router / [post]
func (o *AttackAlarmController) Post() {
	var alarms []map[string]interface{}
	o.UnmarshalJson(&alarms)

	count := 0
	for _, alarm := range alarms {
		alarm["@timestamp"] = time.Now().UnixNano() / 1000000
		err := logs.AddAttackAlarm(alarm)
		if err == nil {
			count++
		}
	}
	o.Serve(map[string]uint64{"count": uint64(count)})
}
