package app

import (
	"admin/server/pkg/e"
	"github.com/gin-gonic/gin"
)

// Gin gin wrapper
type Gin struct {
	C *gin.Context
}

// Response response wrapper
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response wrapper
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})

	return
}
