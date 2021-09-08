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

type CrashController struct {
	controllers.BaseController
}

// @router /search [post]
func (o *CrashController) Search() {
	var param = &logs.SearchCrashParam{}
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
	o.ValidPage(param.Page, param.Perpage)
	if param.Data.StartTime <= 0 {
		o.ServeError(http.StatusBadRequest, "start_time must be greater than 0")
	}
	if param.Data.EndTime <= 0 {
		o.ServeError(http.StatusBadRequest, "end_time must be greater than 0")
	}
	if param.Data.StartTime > param.Data.EndTime {
		o.ServeError(http.StatusBadRequest, "start_time cannot be greater than end_time")
	}
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
		param.Page, param.Perpage, false, logs.CrashAlarmInfo.EsAliasIndex+"-"+param.Data.AppId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to search data from es", err)
	}
	o.Serve(map[string]interface{}{
		"total":      total,
		"total_page": math.Ceil(float64(total) / float64(param.Perpage)),
		"page":       param.Page,
		"perpage":    param.Perpage,
		"data":       result,
	})
}
