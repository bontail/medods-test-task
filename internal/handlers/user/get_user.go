package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetUser(c *gin.Context) {
	token, err := h.GetAccessToken(c)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	user, err := h.UserStorage.GetUser(c, h.UserStorage.WithGUID(token.GUID))
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	h.SendOK(c, user)
}
