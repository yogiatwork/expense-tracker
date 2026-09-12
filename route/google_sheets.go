package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/handler"
)

// WhatsAppRoute defines the routes for WhatsApp related endpoints
func GoogleSheetsRoute(rg *gin.RouterGroup) {
	gsheets := rg.Group("/gsheets")

	gsheets.GET("/metadata", handler.GetSheetMetadataHandler)
	gsheets.POST("/expense", handler.ExpenseEntryHandler)
	gsheets.POST("/income", handler.IncomeEntryHandler)
}
