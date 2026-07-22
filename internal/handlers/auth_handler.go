package handlers

import (
	"net/http"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	result, err := h.service.Login(input)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "username atau password salah"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
