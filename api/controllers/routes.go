package controllers

import (
	"github.com/Musashi-Sakamoto/fullstack/api/middlewares"
	"github.com/gin-gonic/gin"
)

func (s *Server) initializeRoutes() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/", middlewares.SetMiddlewareJSON(), s.Home)
	r.POST("/login", middlewares.SetMiddlewareJSON(), s.Login)

	usersGroup := r.Group("/users", middlewares.SetMiddlewareJSON())
	{
		usersGroup.POST("/", s.CreateUser)
		usersGroup.GET("/", s.GetUsers)
		usersGroup.GET("/:id", s.GetUser)
		usersGroup.PUT("/:id", middlewares.SetMiddlewareAuthentication(), s.UpdateUser)
		usersGroup.DELETE("/:id", middlewares.SetMiddlewareAuthentication(), s.DeleteUser)
	}

	postsGroup := r.Group("/posts", middlewares.SetMiddlewareJSON())
	{
		postsGroup.POST("/", s.CreatePost)
		postsGroup.GET("/", s.GetPosts)
		postsGroup.GET("/:id", s.GetPost)
		postsGroup.PUT("/:id", middlewares.SetMiddlewareAuthentication(), s.UpdatePost)
		postsGroup.DELETE("/:id", middlewares.SetMiddlewareAuthentication(), s.DeletePost)
	}

	r.Run(":8080")
}
