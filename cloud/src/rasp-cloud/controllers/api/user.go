//Copyright 2021-2021 corecna Inc.
package api

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"rasp-cloud/controllers"
	"rasp-cloud/models"
	"strconv"
	"time"
)

type UserController struct {
	controllers.BaseController
}

// @router /login [post]
func (o *UserController) Login() {
	var loginData map[string]string
	o.UnmarshalJson(&loginData)
	logUser := loginData["username"]
	logPasswd := loginData["password"]
	if logUser == "" || logPasswd == "" {
		o.ServeError(http.StatusBadRequest, "username or password cannot be empty")
	}
	if len(logUser) > 512 || len(logPasswd) > 512 {
		o.ServeError(http.StatusBadRequest, "the length of username or password cannot be greater than 512")
	}
	user, err := models.VerifyUser(logUser, logPasswd)
	if err != nil {
		o.ServeError(http.StatusBadRequest, "login failed", err)
	}
	cookie := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(10000))+logUser+"corerasp"+
		strconv.FormatInt(time.Now().UnixNano(), 10))))
	err = models.NewCookie(cookie, user.Id)
	if err != nil {
		o.ServeError(http.StatusUnauthorized, "failed to create cookie", err)
	}
	o.Ctx.SetCookie(models.AuthCookieName, cookie)
	o.ServeWithEmptyData()
}

// @router /default [get,post]
func (o *UserController) CheckDefault() {
	var err error
	var result bool
	var cookie = o.Ctx.GetCookie(models.AuthCookieName)
	if cookie == "" {
		// 解决 token 认证获取不到 cookie 问题，以后多用户此处需要更改
		result, err = models.CheckDefaultPasswordWithDefaultUser()
	} else {
		result, err = models.CheckDefaultPassword(cookie)
	}
	if err != nil {
		o.ServeError(http.StatusBadRequest, "failed to check default password: "+err.Error())
	}
	if result {
		o.Serve(map[string]interface{}{
			"is_default": true,
		})
	} else {
		o.Serve(map[string]interface{}{
			"is_default": false,
		})
	}
}

// @router /islogin [get,post]
func (o *UserController) IsLogin() {
	o.ServeWithEmptyData()
}

// @router /update [post]
func (o *UserController) Update() {
	var param struct {
		OldPwd string `json:"old_password"`
		NewPwd string `json:"new_password"`
	}
	o.UnmarshalJson(&param)
	if param.OldPwd == "" {
		o.ServeError(http.StatusBadRequest, "old_password can not be empty")
	}
	if param.NewPwd == "" {
		o.ServeError(http.StatusBadRequest, "new_password can not be empty")
	}
	err := models.RemoveAllCookie()
	if err != nil {
		o.ServeError(http.StatusBadRequest, err.Error())
	}
	err = models.UpdatePassword(param.OldPwd, param.NewPwd)
	if err != nil {
		o.ServeError(http.StatusBadRequest, err.Error())
	}
	o.ServeWithEmptyData()
}

// @router /logout [get,post]
func (o *UserController) Logout() {
	o.Ctx.SetCookie(models.AuthCookieName, "")
	cookie := o.Ctx.GetCookie(models.AuthCookieName)
	models.RemoveCookie(cookie)
	o.ServeWithEmptyData()
}
