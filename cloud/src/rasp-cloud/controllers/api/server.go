//Copyright 2021-2021 corecna Inc.

package api

import (
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/es"
	"rasp-cloud/models"
	"strings"

	"gopkg.in/mgo.v2"
)

type ServerController struct {
	controllers.BaseController
}

// @router /url/get [post]
func (o *ServerController) GetUrl() {
	serverUrl, err := models.GetServerUrl()
	if err != nil {
		if mgo.ErrNotFound == err {
			o.Serve(models.ServerUrl{AgentUrls: []string{}})
			return
		}
		o.ServeError(http.StatusBadRequest, "failed to get serverUrl", err)
	}
	o.Serve(serverUrl)
}

// @router /url [post]
func (o *ServerController) PutUrl() {
	var serverUrl = &models.ServerUrl{}
	o.UnmarshalJson(&serverUrl)

	if !validHttpUrl(serverUrl.PanelUrl) {
		o.ServeError(http.StatusBadRequest, "Invalid panel url: "+serverUrl.PanelUrl)
	}

	if len(serverUrl.PanelUrl) > 512 {
		o.ServeError(http.StatusBadRequest, "the length of panel url cannot be greater than 512")
	}

	if len(serverUrl.AgentUrls) > 1024 {
		o.ServeError(http.StatusBadRequest, "the count of agent url cannot be greater than 1024")
	}

	if serverUrl.AgentUrls != nil {
		for _, agentUrl := range serverUrl.AgentUrls {
			if len(agentUrl) > 512 {
				o.ServeError(http.StatusBadRequest, "the length of agent url cannot be greater than 512")
			}
			if !validHttpUrl(agentUrl) {
				o.ServeError(http.StatusBadRequest, "Invalid agent url: "+agentUrl)
			}
		}
	}

	err := models.PutServerUrl(serverUrl)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to put server url", err)
	}

	o.Serve(serverUrl)
}

// @router /clear_logs [post]
func (o *ServerController) ClearLogs() {
	docTypeList := []string{"attack-alarm", "report-data", "error-alarm", "policy-alarm", "crash-alarm", "dependency-data"}
	var param struct {
		AppId string `json:"app_id"`
	}

	o.UnmarshalJson(&param)
	if param.AppId == "" {
		o.ServeError(http.StatusBadRequest, "app_id can not be empty")
	}

	for _, docType := range docTypeList {
		index := "real-corerasp-" + docType + "-" + param.AppId
		err := es.DeleteLogs(index)
		if err != nil {
			o.ServeError(http.StatusBadRequest, err.Error())
		}
	}

	o.ServeWithEmptyData()
}

func validHttpUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
