package api

import (
	"admin/core/log"
	"admin/core/rbac"
	"admin/server/pkg/app"
	"admin/server/pkg/e"
	"admin/server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	IDENTITY = "id"
)

type User struct {
	ID       uint
	Username string
	Role     string
}

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

func GetUserInfo(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	userJwt, exist := c.Get(IDENTITY)
	if !exist {
		httpCode = http.StatusBadRequest
		appG.Response(httpCode, errCode, "user jwt claim is empty", nil)
		return
	}

	userClaim := userJwt.(User)
	id := userClaim.ID
	userSrv := service.User{
		ID: id,
	}

	user, err := userSrv.Get()
	if err != nil {
		log.Logger.Error("user", zap.String("err", err.Error()))
		httpCode = http.StatusBadRequest
		appG.Response(httpCode, errCode, "", nil)
		return
	}

	if user.ID == 0 {
		httpCode = http.StatusBadRequest
		appG.Response(httpCode, errCode, "user jwt claim is not correct", nil)
		return
	}

	data := make(map[string]interface{})
	data["introduction"] = user.Remark
	data["name"] = user.Username
	data["displayName"] = user.DisplayName

	var roles []string
	roles = append(roles, user.Role)
	data["roles"] = roles

	routes := rbac.GetRoleRouter(user.Role)
	if len(routes) > 0 {
		data["routes"] = routes
	} else {
		data["routes"] = []struct{}{}
	}

	apis := rbac.GetRoleApi(user.Role)
	if len(apis) > 0 {
		data["routes"] = apis
	} else {
		data["routes"] = []struct{}{}
	}

	appG.Response(httpCode, errCode, "", data)
	return
}

type loginForm struct {
	UserName  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	LoginType string `json:"loginType" validate:"required"`
}

func Login(c *gin.Context) (interface{}, error) {
	var form loginForm

	err := app.BindAndValid(c, &form)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	userSrv := service.User{
		Username:  form.UserName,
		Password:  form.Password,
		LoginType: form.LoginType,
	}

	user, err := userSrv.GetLoginUser()
	if err != nil {
		/* For security */
		log.Logger.Error("user", zap.String("err", err.Error()))
		return nil, fmt.Errorf("%s", "user login failed")
	}

	data := make(map[string]interface{})
	data["id"] = user.ID
	data["username"] = user.Username
	data["role"] = user.Role

	return data, nil
}

func Logout(c *gin.Context) {

}
