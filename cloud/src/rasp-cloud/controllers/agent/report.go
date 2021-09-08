//Copyright 2021-2021 corecna Inc.

package agent

import (
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
)

type ReportController struct {
	controllers.BaseController
}

// @router / [post]
func (o *ReportController) Post() {
	var reportData *models.ReportData
	o.UnmarshalJson(&reportData)
	if reportData.RaspId == "" {
		o.ServeError(http.StatusBadRequest, "rasp_id cannot be empty")
	}
	rasp, err := models.GetRaspById(reportData.RaspId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to get rasp", err)
	}
	if reportData.Time <= 0 {
		o.ServeError(http.StatusBadRequest, "time param must be greater than 0")
	}
	if reportData.RequestSum < 0 {
		o.ServeError(http.StatusBadRequest, "request_sum param cannot be less than 0")
	}
	err = models.AddReportData(reportData, rasp.AppId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to insert report data", err)
	}
	o.Serve(reportData)
}
