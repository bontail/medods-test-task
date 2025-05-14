package handlers

import (
	"encoding/base64"
	"net/netip"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"medods-test-task/internal/models"
)

type SignInData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var data SignInData
	if err := c.ShouldBindJSON(&data); err != nil {
		h.SendFieldErrors(c, err)
		return
	}

	user, err := h.UserStorage.GetUser(c, h.UserStorage.WithUsername(data.Username))
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	if user == nil || !user.ComparePassword(data.Password) {
		h.SendBadRequestError(c, gin.H{"message": "user with this credentials does not exist"})
		return
	}

	addr, err := h.GetAddr(c)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	h.sendTokens(c, user.GUID, addr)
}

func (h *AuthHandler) sendTokens(c *gin.Context, userGUID string, addr netip.Addr) {
	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Duration(h.Cfg.RefreshTokenLifetimeHours) * time.Hour)
	randomValue := uuid.NewString()
	refreshToken, err := models.NewRefreshToken(userGUID, randomValue, createdAt, expiresAt, c.Request.UserAgent(), addr)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	refreshToken.Id, err = h.AuthStorage.InsertToken(c, refreshToken)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	expiresAt = createdAt.Add(time.Duration(h.Cfg.AccessTokenLifetimeSeconds) * time.Second)
	accessToken := models.NewAccessToken(userGUID, refreshToken.Id, expiresAt)
	access, err := accessToken.Encode(h.Cfg.SecretKey)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	h.SendCreated(c, gin.H{
		"access":  access,
		"refresh": base64.StdEncoding.EncodeToString([]byte(randomValue)),
	})
}
