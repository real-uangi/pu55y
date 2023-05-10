// Package api @author uangi 2023-05
package api

import (
	"github.com/real-uangi/pu55y/character"
	"github.com/real-uangi/pu55y/date"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
}

func newResult(code int, message string, data interface{}) Result {
	return Result{
		Code:    code,
		Data:    data,
		Message: message,
		Time:    date.CurrentDateString(),
	}
}

func Success(data interface{}) Result {
	return newResult(200, "success", data)
}

func Fail(err error, message string) Result {
	if err != nil {
		if character.IsNotBlank(message) {
			return newResult(0, message, err.Error())
		}
		return newResult(500, "failed", err.Error())
	}
	if character.IsNotBlank(message) {
		return newResult(0, message, message)
	}
	return newResult(500, "failed", nil)
}

func NotFound(message string) Result {
	if character.IsBlank(message) {
		message = "404 Not Found"
	}
	return newResult(404, message, message)
}
