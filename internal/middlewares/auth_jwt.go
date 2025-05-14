package middlewares

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"medods-test-task/internal/handlers"
	"medods-test-task/internal/models"
	"medods-test-task/internal/server/config"
)

type JWTMiddleware struct {
	handlers.BaseHandler
}

func NewJWTMiddleware(lgr *slog.Logger, cfg *config.Config) JWTMiddleware {
	return JWTMiddleware{
		BaseHandler: handlers.BaseHandler{
			Lgr: lgr,
			Cfg: cfg,
		},
	}
}

func (m JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.Split(c.GetHeader("Authorization"), " ")
		if len(auth) != 2 {
			m.SendAuthorizationError(c, errors.New("Invalid len in Authorization "+strconv.Itoa(len(auth))))
			return
		}
		accessToken, err := models.NewAccessTokenFromString(auth[1], m.Cfg.SecretKey)
		if err != nil {
			m.SendAuthorizationError(c, err)
			return
		}

		c.Set("token", accessToken)
		c.Next()
	}
}
