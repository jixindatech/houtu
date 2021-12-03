package util

import (
	"admin/core/log"
	"admin/server/pkg/app"
	"admin/server/pkg/e"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type user struct {
	ID       uint
	Username string
	Role     string
}

var identityKey = "id" //primary key for gin

func GetJwtMiddleWare(login func(c *gin.Context) (interface{}, error), logout func(c *gin.Context)) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		Unauthorized: func(c *gin.Context, code int, message string) {
			appG := app.Gin{C: c}
			appG.Response(code, e.ERROR, message, nil)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			//data["id"] = id
			//data["username"] = username
			//data["role"] = role
			return login(c)
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"username":  v.Username,
					"role":      v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			appG := app.Gin{C: c}
			data := make(map[string]interface{})
			data["token"] = token
			data["expire"] = expire.Format(time.RFC3339)
			appG.Response(http.StatusOK, e.SUCCESS, "", data)
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &user{
				ID:       claims[identityKey].(uint),
				Username: claims["username"].(string),
				Role:     claims["role"].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// rbac Test here
			if v, ok := data.(*user); ok && v.Username == "admin" {
				return true
			}

			return false
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			appG := app.Gin{C: c}
			data := make(map[string]interface{})
			data["token"] = token
			data["expire"] = expire.Format(time.RFC3339)
			appG.Response(http.StatusOK, e.SUCCESS, "", data)
		},
		LogoutResponse: func(c *gin.Context, code int) {
			appG := app.Gin{C: c}
			logout(c)
			appG.Response(http.StatusOK, e.SUCCESS, "", nil)
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

func NoRoute(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)
	log.Logger.Info("route",
		zap.String("username", username),
		zap.String("method", c.Request.Method),
		zap.String("uri", c.Request.URL.Path))
	appG := app.Gin{C: c}
	appG.Response(http.StatusNotFound, e.ERROR, "", nil)
}
