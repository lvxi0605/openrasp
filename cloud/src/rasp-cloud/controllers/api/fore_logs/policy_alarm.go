//Copyright 2021-2021 corecna Inc.

package fore_logs

import (
	"encoding/json"
	"math"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
	"rasp-cloud/models/logs"
)

// Operations about policy alarm message
type PolicyAlarmController struct {
	controllers.BaseController
}

// @router /search [post]
func (o *PolicyAlarmController) Search() {
	var param = &logs.SearchPolicyParam{}
	o.UnmarshalJson(&param)
	if param.Data == nil {
		o.ServeError(http.StatusBadRequest, "search data can not be empty")
	}
	if param.Data.AppId != "" {
		_, err := models.GetAppById(param.Data.AppId)
		if err != nil {
			o.ServeError(http.StatusBadRequest, "cannot get the app: "+param.Data.AppId, err)
		}
	} else {
		param.Data.AppId = "*"
	}
	if param.Data.StartTime <= 0 {
		o.ServeError(http.StatusBadRequest, "start_time must be greater than 0")
	}
	if param.Data.EndTime <= 0 {
		o.ServeError(http.StatusBadRequest, "end_time must be greater than 0")
	}
	if param.Data.StartTime > param.Data.EndTime {
		o.ServeError(http.StatusBadRequest, "start_time cannot be greater than end_time")
	}
	o.ValidPage(param.Page, param.Perpage)
	content, err := json.Marshal(param.Data)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to encode search data", err)
	}
	var searchData map[string]interface{}
	err = json.Unmarshal(content, &searchData)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to decode search data", err)
	}
	delete(searchData, "start_time")
	delete(searchData, "end_time")
	delete(searchData, "app_id")
	total, result, err := logs.SearchLogs(param.Data.StartTime, param.Data.EndTime, false, searchData, "event_time",
		param.Page, param.Perpage, false, logs.PolicyAlarmInfo.EsAliasIndex+"-"+param.Data.AppId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to search data from es", err)
	}
	// golang禁止循环导包，因此es.go中不能有访问mongo的操作
	// 遍历result，加入rasp_version
	for idx, r := range result {
		if r["rasp_id"] != nil {
			raspId := r["rasp_id"].(string)
			rasp, err := models.GetRaspById(raspId)
			if err == nil {
				result[idx]["rasp_version"] = rasp.Version
			}
		}
	}
	o.Serve(map[string]interface{}{
		"total":      total,
		"total_page": math.Ceil(float64(total) / float64(param.Perpage)),
		"page":       param.Page,
		"perpage":    param.Perpage,
		"data":       result,
	})
}
