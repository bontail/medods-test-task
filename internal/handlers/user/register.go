package handlers

import (
	"github.com/gin-gonic/gin"

	"medods-test-task/internal/models"
)

type RegisterData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var data RegisterData
	if err := c.ShouldBindJSON(&data); err != nil {
		h.SendFieldErrors(c, err)
		return
	}

	user, err := models.NewUser(data.Username, data.Password)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}

	exists, err := h.UserStorage.ExistsUser(c, user.Username)
	if err != nil {
		h.SendInternalError(c, err)
		return
	}
	if exists {
		h.SendBadRequestError(c, gin.H{"username": "already exists"})
		return
	}

	if err = h.UserStorage.InsertUser(c, user); err != nil {
		h.SendInternalError(c, err)
		return
	}

	h.SendCreated(c, gin.H{"Message": "Success register user"})
}
