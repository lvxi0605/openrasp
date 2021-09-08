//Copyright 2021-2021 corecna Inc.

package api

import (
	"io/ioutil"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
	"rasp-cloud/mongo"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Operations about plugin
type PluginController struct {
	controllers.BaseController
}

// @router / [post]
func (o *PluginController) Upload() {
	appId := o.GetString("app_id")
	if appId == "" {
		o.ServeError(http.StatusBadRequest, "app_id can not be empty")
	}
	_, err := models.GetAppById(appId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to get app", err)
	}

	uploadFile, info, err := o.GetFile("plugin")
	if err != nil {
		o.ServeError(http.StatusBadRequest, "parse uploadFile error", err)
	}
	if uploadFile == nil {
		o.ServeError(http.StatusBadRequest, "must have the plugin parameter")
	}
	defer uploadFile.Close()

	if info.Size == 0 {
		o.ServeError(http.StatusBadRequest, "the upload file cannot be empty")
	}
	pluginContent, err := ioutil.ReadAll(uploadFile)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to read upload plugin", err)
	}

	latestPlugin, err := models.AddPlugin(pluginContent, appId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to add plugin", err)
	}
	models.AddOperation(appId, models.OperationTypeUploadPlugin, strings.Split(o.Ctx.Request.RemoteAddr, ":")[0],
		"New plugin uploaded: "+latestPlugin.Id)
	o.Serve(latestPlugin)
}

// @router /get [post]
func (o *PluginController) Get() {
	var param map[string]string
	o.UnmarshalJson(&param)
	pluginId := param["id"]
	if pluginId == "" {
		o.ServeError(http.StatusBadRequest, "plugin_id cannot be empty")
	}
	plugin, err := models.GetPluginById(pluginId, false)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to get plugin", err)
	}
	o.Serve(plugin)
}

// @router /download [get]
func (o *PluginController) Download() {
	pluginId := o.GetString("id")
	plugin, err := models.GetPluginById(pluginId, true)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to get plugin", err)
	}
	o.Ctx.Output.Header("Content-Type", "text/plain")
	if plugin.Name == "" {
		plugin.Name = "plugin"
	}
	o.Ctx.Output.Header("Content-Disposition", "attachment;filename="+plugin.Name+"-"+plugin.Version+".js")
	if len(plugin.OriginContent) != 0 {
		o.Ctx.Output.Body([]byte(plugin.OriginContent))
	} else {
		o.Ctx.Output.Body([]byte(plugin.Content))
	}
}

// @router /algorithm/config [post]
func (o *PluginController) UpdateAppAlgorithmConfig() {
	var param struct {
		PluginId string                 `json:"id"`
		Config   map[string]interface{} `json:"config"`
	}
	o.UnmarshalJson(&param)
	if param.PluginId == "" {
		o.ServeError(http.StatusBadRequest, "plugin id can not be empty")
	}
	if param.Config == nil {
		o.ServeError(http.StatusBadRequest, "config can not be empty")
	}
	appId, err := models.UpdateAlgorithmConfig(param.PluginId, param.Config)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to update algorithm config", err)
	}
	models.AddOperation(appId, models.OperationTypeUpdateAlgorithmConfig,
		strings.Split(o.Ctx.Request.RemoteAddr, ":")[0], "Algorithm config updated for plugin: "+param.PluginId)
	o.ServeWithEmptyData()
}

// @router /algorithm/restore [post]
func (o *PluginController) RestoreAlgorithmConfig() {
	var param map[string]string
	o.UnmarshalJson(&param)
	pluginId := param["id"]
	if pluginId == "" {
		o.ServeError(http.StatusBadRequest, "plugin_id cannot be empty")
	}
	appId, err := models.RestoreDefaultConfiguration(pluginId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to restore the default algorithm config", err)
	}
	models.AddOperation(appId, models.OperationTypeRestorePlugin, strings.Split(o.Ctx.Request.RemoteAddr, ":")[0],
		"Restored algorithm config for plugin: "+pluginId)
	o.ServeWithEmptyData()
}

// @router /delete [post]
func (o *PluginController) Delete() {
	var param map[string]string
	o.UnmarshalJson(&param)
	pluginId := param["id"]
	if pluginId == "" {
		o.ServeError(http.StatusBadRequest, "plugin_id cannot be empty")
	}
	plugin, err := models.GetPluginById(pluginId, false)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "can not get the plugin", err)
	}
	var app *models.App
	err = mongo.FindOne("app", bson.M{"selected_plugin_id": pluginId}, &app)
	if err != nil && err != mgo.ErrNotFound {
		o.ServeError(http.StatusBadRequest, "failed to get app", err)
	}
	if app != nil {
		o.ServeError(http.StatusBadRequest, "Unable to delete a plugin in use. Plugin is used by appid: "+app.Id)
	}
	err = models.DeletePlugin(pluginId)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to delete the plugin", err)
	}
	models.AddOperation(plugin.AppId, models.OperationTypeDeletePlugin, strings.Split(o.Ctx.Request.RemoteAddr, ":")[0],
		"Deleted plugin: "+plugin.Id)
	o.ServeWithEmptyData()
}
