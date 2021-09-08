//Copyright 2021-2021 corecna Inc.

package api

import (
	"math"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
)

type DependencyController struct {
	controllers.BaseController
}

// @router /search [post]
func (o *DependencyController) Search() {
	param := o.handleSearchParam()
	total, result, err := models.SearchDependency(param.Data.AppId, param)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "search dependency failed", err)
	}
	o.Serve(map[string]interface{}{
		"total":      total,
		"total_page": math.Ceil(float64(total) / float64(param.Perpage)),
		"page":       param.Page,
		"perpage":    param.Perpage,
		"data":       result,
	})
}

// @router /aggr [post]
func (o *DependencyController) AggrWithSearch() {
	param := o.handleSearchParam()
	total, result, err := models.AggrDependencyByQuery(param.Data.AppId, param)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "aggregation dependency failed", err)
	}
	o.Serve(map[string]interface{}{
		"total":      total,
		"total_page": math.Ceil(float64(total) / float64(param.Perpage)),
		"page":       param.Page,
		"perpage":    param.Perpage,
		"data":       result,
	})
}

func (o *DependencyController) handleSearchParam() *models.SearchDependencyParam {
	var param models.SearchDependencyParam
	o.UnmarshalJson(&param)
	o.ValidPage(param.Page, param.Perpage)
	o.ValidParam(&param)
	o.ValidParam(param.Data)
	return &param
}
