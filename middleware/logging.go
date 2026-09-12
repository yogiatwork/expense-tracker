package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate request id
		requestId := uuid.New().String()
		c.Set("requestId", requestId)
		c.Header("X-Request-ID", requestId)

		// capture start time for latency calculation
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// process request
		c.Next()

		// request duration
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// log attributes
		attrs := []slog.Attr{
			slog.String("requestId", requestId),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.String("clientIP", c.ClientIP()),
			slog.String("userAgent", c.Request.UserAgent()),
			slog.Int("contentLength", c.Writer.Size()),
		}

		// add err if exists
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.String()))
		}

		// log based on status code
		statusCode := c.Writer.Status()
		msg := "Request completed"
		switch {
		case statusCode >= 500:
			logger.LogAttrs(c.Request.Context(), slog.LevelError, msg, attrs...)
		case statusCode >= 400:
			logger.LogAttrs(c.Request.Context(), slog.LevelWarn, msg, attrs...)
		default:
			logger.LogAttrs(c.Request.Context(), slog.LevelInfo, msg, attrs...)
		}
	}
}
