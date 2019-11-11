package middlewares

import (
	"errors"
	"net/http"

	"github.com/Musashi-Sakamoto/fullstack/api/auth"

	"github.com/gin-gonic/gin"
)

func SetMiddlewareJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func SetMiddlewareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.New("Unauthorized").Error(),
			})
			return
		}
		c.Next()
	}
}
