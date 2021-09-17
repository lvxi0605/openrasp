package logs

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"rasp-cloud/conf"
)

var (
	ErrorAlarmInfo = AlarmLogInfo{
		EsType:       "error-alarm",
		EsIndex:      "corerasp-error-alarm",
		EsAliasIndex: "real-corerasp-error-alarm",
		AlarmBuffer:  make(chan map[string]interface{}, conf.AppConfig.AlarmBufferSize),
		FileLogger:   initAlarmFileLogger("corerasp-logs/error-alarm", "error.log"),
	}
)

func init() {
	registerAlarmInfo(&ErrorAlarmInfo)
}

func AddErrorAlarm(alarm map[string]interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("failed to add error alarm: ", r)
		}
	}()
	idContent := ""
	idContent += fmt.Sprint(alarm["rasp_id"])
	idContent += fmt.Sprint(alarm["error_code"])
	idContent += fmt.Sprint(alarm["message"])
	idContent += fmt.Sprint(alarm["stack_trace"])
	alarm["upsert_id"] = fmt.Sprintf("%x", md5.Sum([]byte(idContent)))
	err := AddLogWithKafka(AttackAlarmInfo.EsType, alarm)
	if err != nil {
		return err
	}
	return AddAlarmFunc(ErrorAlarmInfo.EsType, alarm)
}
