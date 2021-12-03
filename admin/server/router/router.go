package router

import (
	"admin/core/log"
	"admin/server/pkg/app"
	"admin/server/router/api"
	"admin/server/util"
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

	authMiddleware, err := util.GetJwtMiddleWare(nil, nil)
	if err != nil {
		return nil, err
	}
	r.NoRoute(authMiddleware.MiddlewareFunc(), util.NoRoute)

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth := authMiddleware.MiddlewareFunc
	apis := r.Group("/", auth())
	{
		apis.POST("/user", api.AddUser)
	}

	return r, nil
}
