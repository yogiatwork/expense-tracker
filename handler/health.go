package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ready(c *gin.Context) {
	// Simulate a readiness check (e.g., database connection, external service availability)
	// For demonstration purposes, we'll just use a simple time-based check.
	if time.Now().Second()%2 == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	}
}

func Live(c *gin.Context) {
	// Simulate a liveness check (e.g., application health)
	// For demonstration purposes, we'll just return a simple response.
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}
