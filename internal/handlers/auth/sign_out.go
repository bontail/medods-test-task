package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) SignOut(c *gin.Context) {
	token, err := h.GetAccessToken(c)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	err = h.AuthStorage.BlockedAllTokens(c, token.GUID)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	h.SendRequest(c, http.StatusCreated, gin.H{"Message": "You have successfully signed out"})
}
