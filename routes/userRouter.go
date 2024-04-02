package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users")
	incomingRoutes.GET("/users/:user_id")
	incomingRoutes.POST("/users/signUp")
	incomingRoutes.POST("/users/login")
}
