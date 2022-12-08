package routers

import (
	"res/controller"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.Engine) {
	router.GET("/invoices", controller.GetInvoices())
	router.GET("/invoices/:invoice_id", controller.GetInvoice())
	router.POST("/invoices", controller.CreateInvoice())
	router.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
}
