package routes

import "github.com/gin-gonic/gin"

func InvoiceRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices")
	incomingRoutes.GET("/invoice/:invoice_id")
	incomingRoutes.POST("/invoices")
	incomingRoutes.PATCH("/invoices/:invoice_id")
}
