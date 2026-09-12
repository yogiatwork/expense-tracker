package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/model"
	"github.com/yogiatwork/expense-tracker/service"
)

// lets log the request body and headers for debugging purposes
func GetMeHandler(c *gin.Context) {
	id, err := service.GetBotInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"bot_id": id})
}

func RegisterTelegramWebhook(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func ReceiveUpdates(c *gin.Context) {
	// this is telegram webhook handler, it will receive updates from telegram and log them
	var update model.Update
	if err := c.BindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if update.Message.Chat.Type == "private" {
		service.ProcessMsg(update.Message)
	}

	c.JSON(200, gin.H{"status": "ok"})
	slog.Info("received update", slog.Any("update", update))
}
