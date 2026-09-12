package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// get stack trace
			if err := recover(); err != nil {
				// get stack trace
				stack := debug.Stack()

				logger.Error("Panic recovered",
					slog.String("error", fmt.Sprintf("%v", err)),
					slog.String("stack", string(stack)),
					slog.String("requestId", c.GetString("requestId")),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("query", c.Request.URL.RawQuery),
					slog.Int("status", c.Writer.Status()),
					slog.String("clientIP", c.ClientIP()),
					slog.String("userAgent", c.Request.UserAgent()),
				)

				// response based on environment
				respnose := gin.H{"error": "Internal Server Error"}

				if gin.Mode() == gin.DebugMode {
					respnose["stack"] = string(stack)
				}

				c.AbortWithStatusJSON(500, respnose)
			}

		}()

		c.Next()
	}
}
