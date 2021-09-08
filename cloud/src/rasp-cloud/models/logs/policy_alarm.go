//Copyright 2021-2021 corecna Inc.

package logs

import (
	"crypto/md5"
	"fmt"
	"rasp-cloud/conf"

	"github.com/astaxie/beego"
)

var (
	PolicyAlarmInfo = AlarmLogInfo{
		EsType:       "policy-alarm",
		EsIndex:      "openrasp-policy-alarm",
		EsAliasIndex: "real-openrasp-policy-alarm",
		AlarmBuffer:  make(chan map[string]interface{}, conf.AppConfig.AlarmBufferSize),
		FileLogger:   initAlarmFileLogger("openrasp-logs/policy-alarm", "policy.log"),
	}
)

func init() {
	registerAlarmInfo(&PolicyAlarmInfo)
}

func AddPolicyAlarm(alarm map[string]interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("failed to add policy alarm: ", r)
		}
	}()
	putStackMd5(alarm, "policy_params")
	idContent := ""
	idContent += fmt.Sprint(alarm["rasp_id"])
	idContent += fmt.Sprint(alarm["policy_id"])
	idContent += fmt.Sprint(alarm["stack_md5"])
	idContent += fmt.Sprint(alarm["message"])
	if alarm["policy_id"] == "3007" && alarm["policy_params"] != nil {
		if policyParam, ok := alarm["policy_params"].(map[string]interface{}); ok && len(policyParam) > 0 {
			idContent += fmt.Sprint(policyParam["type"])
		}
	} else if alarm["policy_id"] == "3006" && alarm["policy_params"] != nil {
		if policyParam, ok := alarm["policy_params"].(map[string]interface{}); ok && len(policyParam) > 0 {
			idContent += fmt.Sprint(policyParam["connectionString"])
			idContent += fmt.Sprint(policyParam["port"])
			idContent += fmt.Sprint(policyParam["server"])
			idContent += fmt.Sprint(policyParam["hostname"])
			idContent += fmt.Sprint(policyParam["socket"])
			idContent += fmt.Sprint(policyParam["username"])
		}
	} else if alarm["policy_id"] == "3009" && alarm["policy_params"] != nil {
		if policyParam, ok := alarm["policy_params"].(map[string]interface{}); ok && len(policyParam) > 0 {
			idContent += fmt.Sprint(policyParam["webroot"])
		}
	}
	alarm["upsert_id"] = fmt.Sprintf("%x", md5.Sum([]byte(idContent)))
	err := AddLogWithKafka(AttackAlarmInfo.EsType, alarm)
	if err != nil {
		return err
	}
	return AddAlarmFunc(PolicyAlarmInfo.EsType, alarm)
}
