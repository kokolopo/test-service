package http

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine, h *UserHandler) {
	user := r.Group("/users")
	{
		user.POST("/", h.CreateUser)
		user.GET("/", h.GetAllUsers)
		user.GET("/:id", h.GetUserByID)
	}
}
