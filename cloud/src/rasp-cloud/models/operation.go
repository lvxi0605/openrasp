//Copyright 2021-2021 corecna Inc.

package models

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"rasp-cloud/mongo"
	"rasp-cloud/tools"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// operations about app
type Operation struct {
	Id      string `json:"id" bson:"_id,omitempty"`
	TypeId  int    `json:"type_id" bson:"type_id,omitempty"`
	AppId   string `json:"app_id" bson:"app_id,omitempty"`
	Time    int64  `json:"time" bson:"time,omitempty"`
	User    string `json:"user" bson:"user,omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
	Ip      string `json:"ip" bson:"ip,omitempty"`
}

const (
	operationCollectionName = "operation"

	OperationTypeRegisterRasp = 1001 + iota
	OperationTypeDeleteRasp
	OperationTypeRegenerateSecret
	OperationTypeUpdateGenerateConfig
	OperationTypeUpdateWhitelistConfig
	OperationTypeUpdateAlgorithmConfig
	OperationTypeUpdateAlarmConfig
	OperationTypeSetSelectedPlugin
	OperationTypeUploadPlugin
	OperationTypeDeletePlugin
	OperationTypeAddApp
	OperationTypeDeleteApp
	OperationTypeEditApp
	OperationTypeRestorePlugin
	OperationTypeAddStrategy
	OperationTypeDeleteStrategy
	OperationTypeEditStrategy
	OperationTypeSelectStrategy
)

func init() {
	index := &mgo.Index{
		Key:        []string{"app_id"},
		Unique:     false,
		Background: true,
		Name:       "app_id",
	}
	err := mongo.CreateIndex(operationCollectionName, index)
	if err != nil {
		tools.Panic(tools.ErrCodeConfigInitFailed,
			"failed to create app_id index for operation collection", err)
	}

	index = &mgo.Index{
		Key:        []string{"time"},
		Unique:     false,
		Background: true,
		Name:       "time",
	}
	err = mongo.CreateIndex(operationCollectionName, index)
	if err != nil {
		tools.Panic(tools.ErrCodeConfigInitFailed,
			"failed to create operate_time index for operation collection", err)
	}

}

func AddOperation(appId string, typeId int, ip string, content string, userName ...string) error {
	var user string
	var err error
	if len(userName) == 0 {
		user, err = GetLoginUserName()
		if err != nil {
			beego.Error("failed to add operation with content: " + content + ",can not get username: " + err.Error())
			return err
		}
	} else {
		user = userName[0]
	}

	var operation = &Operation{
		AppId:   appId,
		TypeId:  typeId,
		Ip:      ip,
		Id:      generateOperationId(),
		User:    user,
		Time:    time.Now().UnixNano() / 1000000,
		Content: content,
	}
	err = mongo.Insert(operationCollectionName, operation)
	if err != nil {
		beego.Error("failed to add operation with content: " + operation.Content + ",error is: " + err.Error())
	}
	return err
}

func generateOperationId() string {
	random := string(bson.NewObjectId()) + strconv.FormatInt(time.Now().UnixNano(), 10) +
		strconv.Itoa(rand.Intn(10000))
	return fmt.Sprintf("%x", sha1.Sum([]byte(random)))
}

func FindOperation(data *Operation, startTime int64, endTime int64,
	page int, perpage int) (count int, result []Operation, err error) {
	searchData := bson.M{}
	if data.Ip != "" {
		searchData["ip"] = data.Ip
	}
	if data.AppId != "" {
		searchData["app_id"] = data.AppId
	}
	if data.User != "" {
		searchData["user"] = data.User
	}
	if data.TypeId != 0 {
		searchData["type_id"] = data.TypeId
	}
	if data.Id != "" {
		searchData["_id"] = data.Id
	}
	searchData["time"] = bson.M{"$gte": startTime, "$lte": endTime}
	count, err = mongo.FindAll(operationCollectionName, searchData, &result, perpage*(page-1), perpage, "-time")
	if result == nil {
		result = make([]Operation, 0)
	}
	return
}
