package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"example.com/todo-api/internal/todo"
)

func NewRouter(handler *todo.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	todos := router.Group("/todos")
	{
		todos.POST("", handler.Create)
		todos.GET("", handler.List)
		todos.GET("/:id", handler.Get)
		todos.PATCH("/:id", handler.Update)
		todos.DELETE("/:id", handler.Delete)
	}
	return router
}