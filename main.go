package main

import (
	"os"

	_ "github.com/essaubaid/my_first_go_project/config"
	"github.com/essaubaid/my_first_go_project/middleware"
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
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.MenuRouter(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.FoodRouter(router)
	routes.InvoiceRoutes(router)
	routes.TableRoutes(router)

	router.Run(":" + port)
}
