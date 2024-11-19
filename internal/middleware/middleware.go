package middleware

import (
	"crypto/subtle"
	"github.com/aberniando/xenshield/config"
	loggerPkg "github.com/aberniando/xenshield/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func APIKeyAuth() gin.HandlerFunc {
	cfg, _ := config.GetConfig()
	logger := loggerPkg.GetLogger()

	return func(c *gin.Context) {
		apiKey := c.GetHeader("api-key")
		resp := struct {
			Message string `json:"message"`
		}{}

		if apiKey == "" {
			logger.Error("[Middleware] Missing API key in request")
			resp.Message = "Missing API key in request"
			c.JSON(http.StatusForbidden, resp)
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.APIKey)) != 1 {
			logger.Error("[Middleware] Invalid API key provided")
			resp.Message = "Invalid API key"
			c.JSON(http.StatusForbidden, resp)
			c.Abort()
			return
		}

		c.Next()
	}
}
