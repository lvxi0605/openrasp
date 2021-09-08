//Copyright 2021-2021 corecna Inc.

package api

import (
	"math"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
)

type TokenController struct {
	controllers.BaseController
}

// @router /get [post]
func (o *TokenController) Get() {
	var param map[string]int
	o.UnmarshalJson(&param)
	page := param["page"]
	perpage := param["perpage"]
	o.ValidPage(page, perpage)

	total, tokens, err := models.GetAllToken(page, perpage)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to get tokens", err)
	}
	if tokens == nil {
		tokens = make([]*models.Token, 0)
	}
	var result = make(map[string]interface{})
	result["total"] = total
	result["total_page"] = math.Ceil(float64(total) / float64(perpage))
	result["page"] = page
	result["perpage"] = perpage
	result["data"] = tokens
	o.Serve(result)
}

// @router / [post]
func (o *TokenController) Post() {
	var token *models.Token
	o.UnmarshalJson(&token)
	if len(token.Description) > 1024 {
		o.ServeError(http.StatusBadRequest, "the length of the token description must be less than 1024")
	}
	var err error
	if token.Token == "" {
		token, err = models.AddToken(token)
		if err != nil {
			o.ServeError(http.StatusBadRequest, "failed to create token", err)
		}
	} else {
		token, err = models.UpdateToken(token)
		if err != nil {
			o.ServeError(http.StatusBadRequest, "failed to update token", err)
		}
	}
	o.Serve(token)
}

// @router /delete [post]
func (o *TokenController) Delete() {
	var token *models.Token
	o.UnmarshalJson(&token)
	if len(token.Token) == 0 {
		o.ServeError(http.StatusBadRequest, "the token param cannot be empty")
	}
	currentToken := o.Ctx.Input.Header(models.AuthTokenName)
	if currentToken == token.Token {
		o.ServeError(http.StatusBadRequest, "can not delete the token currently in use")
	}
	token, err := models.RemoveToken(token.Token)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to remove token", err)
	}
	o.Serve(token)
}
