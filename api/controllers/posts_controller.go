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

func (server *Server) CreatePost(c *gin.Context) {
	post := models.Post{}
	err := c.BindJSON(&post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}
	if uid != post.AuthorID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(http.StatusText(http.StatusUnauthorized)).Error(),
		})
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": formattedError.Error(),
		})
		return
	}
	c.Writer.Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.URL.Path, postCreated.ID))
	c.JSON(http.StatusCreated, gin.H{
		"postCreated": postCreated,
	})
}

func (server *Server) GetPosts(c *gin.Context) {
	post := models.Post{}
	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (server *Server) GetPost(c *gin.Context) {
	id := c.Param("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"postReceived": postReceived,
	})
}

func (server *Server) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}

	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": errors.New("Post not found").Error(),
		})
		return
	}

	if uid != post.AuthorID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}

	postUpdate := models.Post{}
	err = c.BindJSON(&postUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	if uid != postUpdate.AuthorID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	postUpdate.ID = post.ID
	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": formattedError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"postUpdated": postUpdated,
	})
}

func (server *Server) DeletePost(c *gin.Context) {
	id := c.Param("id")

	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}

	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": errors.New("Not found").Error(),
		})
		return
	}

	if uid != post.AuthorID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Unauthorized").Error(),
		})
		return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Entity", fmt.Sprintf("%d", pid))
	c.JSON(http.StatusNoContent, "")
}
