//Copyright 2021-2021 corecna Inc.

package api

import (
	"math"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
)

type OperationController struct {
	controllers.BaseController
}

// @router /search [post]
func (o *OperationController) Search() {
	var param struct {
		Data      *models.Operation `json:"data"`
		StartTime int64             `json:"start_time"`
		EndTime   int64             `json:"end_time"`
		Page      int               `json:"page"`
		Perpage   int               `json:"perpage"`
	}
	o.UnmarshalJson(&param)
	o.ValidPage(param.Page, param.Perpage)
	if param.Data == nil {
		o.ServeError(http.StatusBadRequest, "search data can not be empty")
	}
	if param.StartTime <= 0 {
		o.ServeError(http.StatusBadRequest, "start_time must be greater than 0")
	}
	if param.EndTime <= 0 {
		o.ServeError(http.StatusBadRequest, "end_time must be greater than 0")
	}
	if param.StartTime > param.EndTime {
		o.ServeError(http.StatusBadRequest, "start_time cannot be greater than end_time")
	}

	var result = make(map[string]interface{})
	total, operations, err := models.FindOperation(param.Data, param.StartTime, param.EndTime, param.Page, param.Perpage)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "Failed to get plugin list", err)
	}
	result["total"] = total
	result["total_page"] = math.Ceil(float64(total) / float64(param.Perpage))
	result["page"] = param.Page
	result["perpage"] = param.Perpage
	result["data"] = operations
	o.Serve(result)
}
