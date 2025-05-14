package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/netip"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"medods-test-task/internal/models"
	"medods-test-task/internal/notificator"
	"medods-test-task/internal/server/config"
)

type BaseHandler struct {
	Lgr *slog.Logger
	Cfg *config.Config
	Ntf notificator.Notificator
}

func (h *BaseHandler) SendRequest(c *gin.Context, code int, obj interface{}) {
	c.AbortWithStatusJSON(code, obj)
}

func (h *BaseHandler) SendOK(c *gin.Context, obj interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, obj)
}

func (h *BaseHandler) SendCreated(c *gin.Context, obj interface{}) {
	c.AbortWithStatusJSON(http.StatusCreated, obj)
}

func (h *BaseHandler) SendBadRequestError(c *gin.Context, obj interface{}) {
	h.SendRequest(c, http.StatusBadRequest, obj)
}

func (h *BaseHandler) SendFieldErrors(c *gin.Context, err error) {
	var res any
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		data := make(map[string]string, len(validateErrs))
		for _, e := range validateErrs {
			data[e.Field()] = e.Tag()
		}
		res = data
	} else {
		res = gin.H{"error": err.Error()}
	}

	h.SendBadRequestError(c, res)
}

func (h *BaseHandler) SendInternalError(c *gin.Context, err error) {
	h.SendRequest(c, http.StatusInternalServerError, gin.H{"message": "Internal error"})
	h.Lgr.Error(err.Error())
}

func (h *BaseHandler) SendAuthorizationError(c *gin.Context, err error) {
	h.SendRequest(c, http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	h.Lgr.Error(err.Error())
}

func (h *BaseHandler) SendForbiddenError(c *gin.Context, reason string) {
	h.SendRequest(c, http.StatusForbidden, gin.H{
		"message": "Forbidden", "reason": reason,
	})
}

func (h *BaseHandler) GetAccessToken(c *gin.Context) (*models.AccessToken, error) {
	t, exists := c.Get("token")
	if !exists {
		return nil, errors.New("no token found")
	}

	token, ok := t.(*models.AccessToken)
	if !ok {
		return nil, errors.New("convert access token error")
	}

	return token, nil
}

func (h *BaseHandler) GetAddr(c *gin.Context) (netip.Addr, error) {
	return netip.ParseAddr(c.ClientIP())
}
