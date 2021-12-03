package api

import (
	"admin/core/log"
	"admin/server/pkg/app"
	"admin/server/pkg/e"
	"admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type userForm struct {
	UserName    string `json:"username" validate:"required"`
	DisplayName string `json:"displayName"`
	LoginType   string `json:"loginType" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,phone"`
	Status      int    `json:"status" validate:"required,gte=0,lte=1"`
	Role        string `json:"role" validate:"required"`
	Remark      string `json:"remark"`
}

func AddUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     userForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	err := app.BindAndValid(c, &form)
	if err != nil {
		httpCode = http.StatusBadRequest
		appG.Response(httpCode, e.ERROR, err.Error(), nil)
		return
	}

	userSrv := service.User{
		Username:    form.UserName,
		DisplayName: form.DisplayName,
		LoginType:   form.DisplayName,
		Email:       form.Email,
		Phone:       form.Phone,
		Status:      form.Status,
		Role:        form.Role,
		Remark:      form.Remark,
	}
	err = userSrv.Save()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserAddFailed
		log.Logger.Error("user", zap.String("err", err.Error()))
	}

	appG.Response(httpCode, errCode, "", nil)
}

func Login(c *gin.Context) (interface{}, error) {
	return nil, nil
}
