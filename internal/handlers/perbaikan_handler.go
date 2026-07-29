package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"invela-be/internal/middlewares"
	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PerbaikanHandler struct {
	service *services.PerbaikanService
}

func NewPerbaikanHandler(service *services.PerbaikanService) *PerbaikanHandler {
	return &PerbaikanHandler{service: service}
}

func (h *PerbaikanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *PerbaikanHandler) Create(c *gin.Context) {
	userIDRaw, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	
	var userID uint
	switch v := userIDRaw.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user id in token"})
		return
	}

	var input services.CreatePerbaikanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, services.ErrPerbaikanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrPerbaikanKerusakanNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_kerusakan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create perbaikan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *PerbaikanHandler) Get(c *gin.Context) {
	id, ok := parsePerbaikanID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrPerbaikanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perbaikan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PerbaikanHandler) Update(c *gin.Context) {
	id, ok := parsePerbaikanID(c)
	if !ok {
		return
	}

	var input services.UpdatePerbaikanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrPerbaikanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perbaikan not found"})
			return
		}

		if errors.Is(err, services.ErrPerbaikanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PerbaikanHandler) Delete(c *gin.Context) {
	id, ok := parsePerbaikanID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrPerbaikanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perbaikan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "perbaikan deleted"})
}

func parsePerbaikanID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
