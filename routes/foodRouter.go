package routes

import (
	"github.com/essaubaid/my_first_go_project/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/food/:food_id", controllers.GetFood())
	incomingRoutes.POST("/food", controllers.CreateFood())
	incomingRoutes.PATCH("/food/:food_id", controllers.UpdateFood())
}
