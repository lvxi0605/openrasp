//Copyright 2021-2021 corecna Inc.
package mongo

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"rasp-cloud/conf"
	"rasp-cloud/tools"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	minMongoVersion = "3.6.0"
	session         *mgo.Session
	DbName          = conf.AppConfig.MongoDBName
)

func init() {
	var err error
	dialInfo := &mgo.DialInfo{
		Addrs:     conf.AppConfig.MongoDBAddr,
		Username:  conf.AppConfig.MongoDBUser,
		Password:  conf.AppConfig.MongoDBPwd,
		Direct:    false,
		Timeout:   time.Second * 20,
		FailFast:  true,
		PoolLimit: conf.AppConfig.MongoDBPoolLimit,
		Database:  DbName,
	}
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		tools.Panic(tools.ErrCodeMongoInitFailed, "Unable to connect to MongoDB server", err)
	}
	info, err := session.BuildInfo()
	if err != nil {
		tools.Panic(tools.ErrCodeMongoInitFailed, "Failed to get MongoDB version", err)
	}
	beego.Info("MongoDB version: " + info.Version)
	if strings.Compare(info.Version, minMongoVersion) < 0 {
		tools.Panic(tools.ErrCodeMongoInitFailed, "MongoDB version lower than "+
			minMongoVersion+" is not supported, current version is "+info.Version, nil)
	}
	if err != nil {
		tools.Panic(tools.ErrCodeMongoInitFailed, "init MongoDB failed", err)
	}

	session.SetMode(mgo.Strong, true)
}

func NewSession() *mgo.Session {
	return session.Copy()
}

func Count(collection string) (int, error) {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).Count()
}

func CreateIndex(collection string, index *mgo.Index) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).EnsureIndex(*index)
}

func Insert(collection string, doc interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).Insert(doc)
}

func UpsertId(collection string, id interface{}, doc interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	_, err := newSession.DB(DbName).C(collection).UpsertId(id, doc)
	return err
}

func FindAll(collection string, query interface{}, result interface{}, skip int, limit int,
	sortFields ...string) (count int, err error) {
	newSession := NewSession()
	defer newSession.Close()
	count, err = newSession.DB(DbName).C(collection).Find(query).Count()
	if err != nil {
		return
	}
	err = newSession.DB(DbName).C(collection).Find(query).Skip(skip).Limit(limit).Sort(sortFields...).All(result)
	return
}

func FindAllWithoutLimit(collection string, query interface{}, result interface{},
	sortFields ...string) (count int, err error) {
	newSession := NewSession()
	defer newSession.Close()
	count, err = newSession.DB(DbName).C(collection).Find(query).Count()
	if err != nil {
		return
	}
	err = newSession.DB(DbName).C(collection).Find(query).Sort(sortFields...).All(result)
	return
}

func FindAllWithSelect(collection string, query interface{}, result interface{}, selector interface{},
	skip int, limit int) (count int, err error) {
	newSession := NewSession()
	defer newSession.Close()
	count, err = newSession.DB(DbName).C(collection).Find(query).Count()
	if err != nil {
		return
	}
	err = newSession.DB(DbName).C(collection).Find(query).Select(selector).Skip(skip).Limit(limit).All(result)
	return
}

func FindSelectWithAggregation(collection string, query interface{}, result interface{}) (err error) {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).Pipe(query).All(result)
}

func FindId(collection string, id string, result interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).FindId(id).One(result)
}

func FindOne(collection string, query interface{}, result interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).Find(query).One(result)
}

func FindAllBySort(collection string, query interface{}, skip int, limit int, result interface{},
	sortFields ...string) (count int, err error) {
	newSession := NewSession()
	defer newSession.Close()
	count, err = newSession.DB(DbName).C(collection).Find(query).Count()
	if err != nil {
		return
	}
	return count, newSession.DB(DbName).C(collection).Find(query).Sort(sortFields...).Skip(skip).Limit(limit).All(result)
}

func UpdateId(collection string, id interface{}, doc interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).UpdateId(id, bson.M{"$set": doc})
}

func RemoveId(collection string, id interface{}) error {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).RemoveId(id)
}

func RemoveAll(collection string, selector interface{}) (*mgo.ChangeInfo, error) {
	newSession := NewSession()
	defer newSession.Close()
	return newSession.DB(DbName).C(collection).RemoveAll(selector)
}

func GenerateObjectId() string {
	random := string(bson.NewObjectId()) +
		strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(10000))
	return fmt.Sprintf("%x", sha1.Sum([]byte(random)))
}
