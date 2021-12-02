package router

import "github.com/gin-gonic/gin"

func Setup(mode string) (g *gin.Engine, err error) {
	r := gin.New()
	gin.SetMode(mode)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r, nil
}
