package main

import (
	"os"

	_ "github.com/essaubaid/my_first_go_project/config"
	"github.com/essaubaid/my_first_go_project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.MenuRouter(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	// router.Use(middleware.Authentication())

	router.Run(":" + port)
}
