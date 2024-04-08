package routes

import (
	"github.com/essaubaid/my_first_go_project/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoice/:invoice_id", controllers.GetInvoice())
	incomingRoutes.POST("/invoice", controllers.CreateInvoice())
	incomingRoutes.PATCH("/invoice/:invoice_id", controllers.UpdateInvoice())
}
