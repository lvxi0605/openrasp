//Copyright 2021-2021 corecna Inc.
package models

import "rasp-cloud/mongo"

const (
	serverUrlId               = "0"
	serverUrlCollectionName   = "server_url"
	serverAgentCollectionName = "server_agent"
)

type ServerUrl struct {
	PanelUrl  string   `json:"panel_url" bson:"panel_url"`
	AgentUrls []string `json:"agent_urls" bson:"agent_urls"`
}

func GetServerUrl() (serverUrl *ServerUrl, err error) {
	err = mongo.FindId(serverUrlCollectionName, serverUrlId, &serverUrl)
	if err == nil && serverUrl.AgentUrls == nil {
		serverUrl.AgentUrls = []string{}
	}
	return
}

func PutServerUrl(serverUrl *ServerUrl) error {
	return mongo.UpsertId(serverUrlCollectionName, serverUrlId, &serverUrl)
}
