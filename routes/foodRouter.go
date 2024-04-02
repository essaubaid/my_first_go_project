package routes

import "github.com/gin-gonic/gin"

func FoodRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods")
	incomingRoutes.GET("/food/:food_id")
	incomingRoutes.POST("/foods")
	incomingRoutes.PATCH("/foods/:food_id")
}
