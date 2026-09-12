package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/config"
)

// it just receives message date and prints to logs and returns 200 OK
func WebhookWhatsAppMessage(c *gin.Context) {
	var message map[string]interface{}
	if err := c.BindJSON(&message); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Log the received message
	slog.Info("Received WhatsApp message", "message", message)

	// Respond with 200 OK
	c.JSON(200, gin.H{"status": "received"})
}

func VerifyWhatsAppWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	// Verify the mode and token
	if mode == "subscribe" && token == config.AppConfig.WhatsappApi.VerifyToken {
		slog.Info("Webhook verified successfully")
		c.String(200, challenge)
	} else {
		slog.Warn("Webhook verification failed")
		c.String(403, "Forbidden")
	}
}
