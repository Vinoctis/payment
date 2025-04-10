package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ResCode int64

const (
  CodeSuccess  ResCode = 1000
  CodeInvalidParam  ResCode = 1001
  CodeUserExist  ResCode = 1002
  CodeUserNotExist  ResCode = 1003
  CodeInvalidPassword   ResCode = 1004
  CodeServerBusy   ResCode = 1005
  CodeNeedLogin   ResCode = 1006
  CodeInvalidToken   ResCode = 1007
)

var codeMsgMap = map[ResCode]string{
  CodeSuccess:         "success",
  CodeInvalidParam:    "请求参数错误",
  CodeUserExist:       "用户名已存在",
  CodeUserNotExist:    "用户名不存在",
  CodeInvalidPassword: "用户名或密码错误",
  CodeServerBusy:      "服务繁忙",
  CodeNeedLogin:       "需要登录",
  CodeInvalidToken:    "无效的token",
}

func (c ResCode) Msg() string {
  msg, ok := codeMsgMap[c]
  if !ok {
    msg = codeMsgMap[CodeServerBusy]
  }
  return msg
}

type ResponseData struct {
	Code ResCode `json:"code"`
	Msg interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg: code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg: msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg: CodeSuccess.Msg(),
		Data: data,
	})
}