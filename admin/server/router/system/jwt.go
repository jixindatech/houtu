package system

import (
	"admin/core/log"
	"admin/core/rbac"
	"admin/server/pkg/app"
	"admin/server/pkg/e"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func GetJwtMiddleWare(login func(c *gin.Context) (interface{}, error), logout func(c *gin.Context)) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IDENTITY,
		Unauthorized: func(c *gin.Context, code int, message string) {
			appG := app.Gin{C: c}
			appG.Response(code, e.ERROR, message, nil)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return login(c)
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*jwtUser); ok {
				return jwt.MapClaims{
					IDENTITY:   v.ID,
					"username": v.Username,
					"role":     v.Role,
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
			return &jwtUser{
				ID:       uint(claims[IDENTITY].(float64)),
				Username: claims["username"].(string),
				Role:     claims["role"].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			enforcer := rbac.GetEnforcer()
			if enforcer == nil {
				return false
			}

			user := data.(*jwtUser)
			sub := user.Role
			uri := c.Request.URL.Path
			method := c.Request.Method

			ok, err := enforcer.Enforce(sub, uri, method)
			if ok {
				return true
			}

			if err != nil {
				log.Logger.Error("jwt", zap.String("enforce", err.Error()))
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

		TokenLookup:   "header: X-Token",
		TokenHeadName: "token",
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
