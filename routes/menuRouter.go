package routes

import "github.com/gin-gonic/gin"

func MenuRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus")
	incomingRoutes.GET("/menu/:menu_id")
	incomingRoutes.POST("/menus")
	incomingRoutes.PATCH("/menus/:menu_id")
}
