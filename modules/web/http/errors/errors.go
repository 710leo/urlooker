package errors

import (
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// 401
func NotLoginError(msg ...string) Error {
	return _build(http.StatusUnauthorized, "unauthorized", msg...)
}

// 400
func BadRequestError(msg ...string) Error {
	return _build(http.StatusBadRequest, "bad request", msg...)
}

// 403
func NoPrivError(msg ...string) Error {
	return _build(http.StatusForbidden, "forbidden", msg...)
}

// 500
func InternalServerError(msg ...string) Error {
	return _build(http.StatusInternalServerError, "internal server error", msg...)
}

func _build(code int, defval string, custom ...string) Error {
	msg := defval
	if len(custom) > 0 {
		msg = custom[0]
	}
	return Error{
		Code: code,
		Msg:  msg,
	}
}

func MaybePanic(err error) {
	_, whichFile, line, _ := runtime.Caller(1)
	arr := strings.Split(whichFile, "/")
	file := arr[len(arr)-1]
	t := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	if err != nil {
		panic(Error{Msg: err.Error(), Time: t, File: file, Line: line})
	}
}

func Panic(msg string) {
	_, whichFile, line, _ := runtime.Caller(1)
	arr := strings.Split(whichFile, "/")
	file := arr[len(arr)-1]
	t := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	panic(Error{Msg: msg, Time: t, File: file, Line: line})
}
