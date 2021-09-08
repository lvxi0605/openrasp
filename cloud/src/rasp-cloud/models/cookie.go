//Copyright 2021-2021 corecna Inc.

package models

import (
	"rasp-cloud/conf"
	"rasp-cloud/mongo"
	"rasp-cloud/tools"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Cookie struct {
	Id     string    `json:"id" bson:"_id"`
	UserId string    `json:"user_id" bson:"user_id"`
	Time   time.Time `json:"time" bson:"time"`
}

const (
	cookieCollectionName = "cookie"
	AuthCookieName       = "RASP_AUTH_ID"
)

func init() {
	index := &mgo.Index{
		Key:         []string{"time"},
		Background:  true,
		Name:        "time",
		ExpireAfter: time.Duration(conf.AppConfig.CookieLifeTime) * time.Hour,
	}
	err := mongo.CreateIndex(cookieCollectionName, index)
	if err != nil {
		tools.Panic(tools.ErrCodeMongoInitFailed, "failed to create index for app collection", err)
	}
}

func NewCookie(id string, userId string) error {
	return mongo.Insert(cookieCollectionName, &Cookie{Id: id, UserId: userId, Time: time.Now()})
}

func HasCookie(id string) (bool, error) {
	var result *Cookie
	err := mongo.FindId(cookieCollectionName, id, &result)
	if err != nil || result == nil {
		return false, err
	}
	return true, err
}

func RemoveCookie(id string) error {
	return mongo.RemoveId(cookieCollectionName, id)
}

func RemoveAllCookie() error {
	_, err := mongo.RemoveAll(cookieCollectionName, bson.M{})
	return err
}
