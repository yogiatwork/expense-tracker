package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/handler"
)

func WhatsAppRoute(r *gin.RouterGroup) {
	whatsapp := r.Group("/whatsapp")

	whatsapp.GET("/message", handler.VerifyWhatsAppWebhook)
	whatsapp.POST("/message", handler.WebhookWhatsAppMessage)
}
