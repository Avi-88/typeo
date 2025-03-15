package main

import (
	"fmt"
	"typeo/config"
	"typeo/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	// config.ConnectDB()
	fmt.Println("Server is up and running on port 8080")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "Hello")
	})
	r.GET("/type", config.LoadSocketServer)
	routes.SetUpAuthRoutes(r)
	routes.SetUpUserRoutes(r)
	r.Run(":8080")
}
