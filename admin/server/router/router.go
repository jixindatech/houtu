package router

import (
	"admin/core/log"
	"admin/server/pkg/app"
	"admin/server/router/system"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func Setup(mode string) (g *gin.Engine, err error) {
	err = app.SetupValidate()
	if err != nil {
		return nil, err
	}

	r := gin.New()
	gin.SetMode(mode)
	r.Use(ginzap.Ginzap(log.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log.Logger, true))

	authMiddleware, err := system.GetJwtMiddleWare(system.Login, system.Logout)
	if err != nil {
		return nil, err
	}
	r.NoRoute(authMiddleware.MiddlewareFunc(), system.NoRoute)

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth := authMiddleware.MiddlewareFunc
	systemApis := r.Group("/system", auth())
	{
		systemApis.GET("/user/refresh_token", authMiddleware.RefreshHandler)
		systemApis.POST("/user/logout", authMiddleware.LogoutHandler)
		systemApis.GET("/user/info", system.GetUserInfo)

		systemApis.POST("/user", system.AddUser)
		systemApis.GET("/user", system.GetUsers)
		systemApis.GET("/user/:id", system.GetUser)
		systemApis.PUT("/user/:id", system.UpdateUser)
		systemApis.PUT("/user", system.UpdateUserInfo)
		systemApis.PUT("/user/password/:id", system.UpdateUserPassword)
		systemApis.DELETE("/user/:id", system.DeleteUser)

		systemApis.GET("/email", system.GetEmail)
		systemApis.POST("/email", system.AddEmail)
		systemApis.PUT("/email/:id", system.UpdateEmail)

		systemApis.GET("/ldap", system.GetLdap)
		systemApis.POST("/ldap", system.AddLdap)
		systemApis.PUT("/ldap/:id", system.UpdateLdap)

		apsystemApisis.GET("/txsms", system.GetTxsms)
		systemApis.POST("/txsms", system.AddTxsms)
		systemApis.PUT("/txsms/:id", system.UpdateTxsms)

		systemApis.POST("/msg", system.AddMsg)
		systemApis.GET("/msg", system.GetMsgs)
		systemApis.GET("/msg/:id", system.GetMsg)
		systemApis.PUT("/msg/:id", system.UpdateMsg)
		systemApis.DELETE("/msg/:id", system.DeleteMsg)
		systemApis.POST("/msg/:id/user", system.SendMsg)
	}

	return r, nil
}
