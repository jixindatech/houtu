package api

import (
	"admin/server/pkg/app"
	"admin/server/pkg/e"
	"github.com/gin-gonic/gin"
	//"github.com/go-playground/validator/v10"
	"net/http"
)

type userForm struct {
	UserName    string `json:"username" binding:"required"`
	DisplayName string `json:"displayName"`
	LoginType   string `json:"loginType" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"phone"`
	Status      int    `json:"status" binding:"gte=0,lte=1"`
	Role        string `json:"role" binding:"required"`
	Remark      string `json:"remark"`
}

func AddUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		//form     userForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	appG.Response(httpCode, errCode, nil)
}
