package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshData struct {
	Refresh string `json:"refresh" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var data RefreshData
	if err := c.ShouldBindJSON(&data); err != nil {
		h.SendFieldErrors(c, err)
		return
	}

	accessToken, err := h.GetAccessToken(c)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	tokenId, err := accessToken.IntID()
	if err != nil {
		h.SendInternalError(c, err)
	}
	refreshToken, err := h.AuthStorage.GetToken(c, accessToken.GUID, tokenId)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	if refreshToken == nil || refreshToken.CompareSecretValue(data.Refresh) {
		h.SendRequest(c, http.StatusForbidden, gin.H{
			"message": "Forbidden", "reason": "Invalid data or access token pair",
		})
		return
	}
	if refreshToken.UserAgent != c.Request.UserAgent() {
		if err = h.AuthStorage.BlockedAllTokens(c, accessToken.GUID); err != nil {
			h.SendInternalError(c, err)
			return
		}
		h.SendRequest(c, http.StatusForbidden, gin.H{
			"message": "Forbidden", "reason": "New user agent",
		})
		return
	}

	addr, err := h.GetAddr(c)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	if refreshToken.IP != addr {
		go h.Ntf.NewIp(refreshToken.IP, addr, accessToken.GUID)
	}

	if err = h.AuthStorage.BlockedToken(c, refreshToken.Id); err != nil {
		h.SendInternalError(c, err)
		return
	}

	h.sendTokens(c, accessToken.GUID, addr)
}
