package routes

import (
	"github.com/essaubaid/my_first_go_project/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/order", controllers.CreateOrder())
	incomingRoutes.PATCH("/order/:order_id", controllers.UpdateOrder())
}
