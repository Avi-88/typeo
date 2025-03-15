package routes

import (
	"typeo/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpAuthRoutes (router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", controllers.HandleLogin)
		authGroup.POST("/register", controllers.HandleRegistration)
	}
}