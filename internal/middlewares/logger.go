package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"medods-test-task/internal/server/config"
)

type LoggerMiddleware struct {
	lgr *slog.Logger
	cfg *config.Config
}

func NewLoggerMiddleware(lgr *slog.Logger, cfg *config.Config) LoggerMiddleware {
	return LoggerMiddleware{
		lgr: lgr,
		cfg: cfg,
	}
}

func (m LoggerMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		m.lgr.Info(
			"send request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("code", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.String("ua", c.Request.UserAgent()),
		)
	}
}
