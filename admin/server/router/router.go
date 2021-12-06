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
	apis := r.Group("/system", auth())
	{
		apis.GET("/user/refresh_token", authMiddleware.RefreshHandler)
		apis.POST("/user/logout", authMiddleware.LogoutHandler)
		apis.GET("/user/info", system.GetUserInfo)

		apis.POST("/user", system.AddUser)
		apis.GET("/user", system.GetUsers)
		apis.GET("/user/:id", system.GetUser)
		apis.PUT("/user/:id", system.UpdateUser)
		apis.PUT("/user/password/:id", system.UpdateUserPassword)
		apis.DELETE("/user/:id", system.DeleteUser)

		apis.GET("/email", system.GetEmail)
		apis.POST("/email", system.AddEmail)
		apis.PUT("/email/:id", system.UpdateEmail)
	}

	return r, nil
}
