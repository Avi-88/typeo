package routes

import (
	"typeo/controllers"
	"typeo/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUpUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users") 
	userGroup.Use(middlewares.VerifyToken)
	{
		userGroup.GET("all", controllers.GetUsers)
	}
}