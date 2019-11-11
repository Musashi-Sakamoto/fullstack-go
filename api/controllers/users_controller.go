package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Musashi-Sakamoto/fullstack/api/auth"
	"github.com/Musashi-Sakamoto/fullstack/api/models"
	"github.com/Musashi-Sakamoto/fullstack/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreateUser(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": formattedError.Error(),
		})
		return
	}
	c.Writer.Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, userCreated.ID))
	c.JSON(http.StatusCreated, gin.H{
		"userCreated": userCreated,
	})
}

func (server *Server) GetUsers(c *gin.Context) {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (server *Server) GetUser(c *gin.Context) {
	id := c.Param("id")

	uid, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userGotten": userGotten,
	})
}

func (server *Server) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{}
	err = c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}
	if tokenID != uint32(uid) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(http.StatusText(http.StatusUnauthorized)).Error(),
		})
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": formattedError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"updatedUser": updatedUser,
	})
}

func (server *Server) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	uid, err := strconv.ParseUint(id, 10, 64)
	user := models.User{}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(http.StatusText(http.StatusUnauthorized)),
		})
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Entity", fmt.Sprintf("%d", uid))
	c.JSON(http.StatusNoContent, "")
}
