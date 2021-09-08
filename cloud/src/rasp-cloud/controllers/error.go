//Copyright 2021-2021 corecna Inc.

package controllers

type ErrorController struct {
	BaseController
}

func (o *ErrorController) Error404() {
	o.errorStatus(404)
}

func (o *ErrorController) Error500() {
	o.errorStatus(500)
}

func (o *ErrorController) Error503() {
	o.errorStatus(503)
}

func (o *ErrorController) Error502() {
	o.errorStatus(502)
}

func (o *ErrorController) errorStatus(code int) {
	o.ServeStatusCode(code, code)
}
