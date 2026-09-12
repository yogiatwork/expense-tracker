package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/handler"
)

func HealthzRoute(r *gin.RouterGroup) {
	health := r.Group("/health")

	health.GET("/ready", handler.Ready)
	health.GET("/live", handler.Live)
}
