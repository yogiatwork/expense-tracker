package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/handler"
)

func TelegramRoute(rg *gin.RouterGroup) {
	// add more routes here
	telegram := rg.Group("/telegram")

	telegram.GET("/bot", handler.GetMeHandler)
	telegram.GET("/messages", handler.RegisterTelegramWebhook)
	telegram.POST("/messages", handler.ReceiveUpdates)
}
