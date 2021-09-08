//Copyright 2021-2021 corecna Inc.

package models

import (
	"rasp-cloud/mongo"
)

type Token struct {
	Token       string `json:"token" bson:"_id"`
	Description string `json:"description" bson:"description"`
}

const (
	tokenCollectionName = "token"
	AuthTokenName       = "X-OpenRASP-Token"
)

func GetAllToken(page int, perpage int) (count int, result []*Token, err error) {
	count, err = mongo.FindAll(tokenCollectionName, nil, &result, perpage*(page-1), perpage)
	return
}

func HasToken(token string) (bool, error) {
	var result *Token
	err := mongo.FindId(tokenCollectionName, token, &result)
	if err != nil || result == nil {
		return false, err
	}
	return true, err
}

func AddToken(token *Token) (result *Token, err error) {
	token.Token = generateOperationId()
	err = mongo.Insert(tokenCollectionName, token)
	if err != nil {
		return
	}
	return token, err
}

func UpdateToken(token *Token) (result *Token, err error) {
	err = mongo.UpdateId(tokenCollectionName, token.Token, token)
	if err != nil {
		return
	}
	return token, err
}

func RemoveToken(tokenId string) (token *Token, err error) {
	err = mongo.FindId(tokenCollectionName, tokenId, &token)
	if err != nil {
		return
	}
	return token, mongo.RemoveId(tokenCollectionName, tokenId)
}
