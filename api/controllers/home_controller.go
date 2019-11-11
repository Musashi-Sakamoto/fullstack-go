package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome To This Awesome API",
	})
}
