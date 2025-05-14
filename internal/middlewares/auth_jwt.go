package middlewares

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"medods-test-task/internal/models"
	"medods-test-task/internal/server/config"
)

type JWTMiddleware struct {
	lgr *slog.Logger
	cfg *config.Config
}

func NewJWTMiddleware(lgr *slog.Logger, cfg *config.Config) JWTMiddleware {
	return JWTMiddleware{
		lgr: lgr,
		cfg: cfg,
	}
}

func (m JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.Split(c.GetHeader("Authorization"), " ")
		if len(auth) != 2 {
			m.lgr.Info("Authorization Error",
				slog.String("err", "Invalid len in Authorization "+strconv.Itoa(len(auth))),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		accessToken, err := models.NewAccessTokenFromString(auth[1], m.cfg.SecretKey)
		if err != nil {
			m.lgr.Info("Authorization Error",
				slog.String("err", err.Error()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		c.Set("token", accessToken)
		c.Next()
	}
}
