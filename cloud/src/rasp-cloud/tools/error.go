//Copyright 2021-2021 corecna Inc.

package tools

import (
	"os"
	"strconv"

	"github.com/astaxie/beego/logs"
)

const (
	// init error
	ErrCodeLogInitFailed = 30001 + iota
	ErrCodeMongoInitFailed
	ErrCodeESInitFailed
	ErrCodeConfigInitFailed
	ErrCodeStartTypeNotSupport
	ErrCodeGeneratePasswdFailed
	ErrCodeGeoipInit
	ErrCodeResetUserFailed
	ErrCodeInitDefaultAppFailed
	ErrCodeInitChildProcessFailed
	ErrCodeChDirFailed
	ErrCodeGetPidFailed
)

const (
	// api error
	ErrRaspNotFound = 4001 + iota
)

func Panic(errCode int, message string, err error) {
	message = "[" + strconv.Itoa(errCode) + "] " + message
	if err != nil {
		message = message + ": " + err.Error()
	}
	logs.Error(message)
	os.Exit(errCode)
}
