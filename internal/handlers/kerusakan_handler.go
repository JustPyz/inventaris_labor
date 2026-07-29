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

type KerusakanHandler struct {
	service *services.KerusakanService
}

func NewKerusakanHandler(service *services.KerusakanService) *KerusakanHandler {
	return &KerusakanHandler{service: service}
}

func (h *KerusakanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kerusakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *KerusakanHandler) Create(c *gin.Context) {
	userIDRaw, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	
	// Convert based on how it's stored in JWT claims
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

	var input services.CreateKerusakanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, services.ErrKerusakanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrKerusakanItemInstanceNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_item_instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create kerusakan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *KerusakanHandler) Get(c *gin.Context) {
	id, ok := parseKerusakanID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrKerusakanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kerusakan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kerusakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *KerusakanHandler) Update(c *gin.Context) {
	id, ok := parseKerusakanID(c)
	if !ok {
		return
	}

	var input services.UpdateKerusakanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrKerusakanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kerusakan not found"})
			return
		}

		if errors.Is(err, services.ErrKerusakanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrKerusakanItemInstanceNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_item_instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update kerusakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *KerusakanHandler) Delete(c *gin.Context) {
	id, ok := parseKerusakanID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrKerusakanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kerusakan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete kerusakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kerusakan deleted"})
}

func parseKerusakanID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}

func (h *KerusakanHandler) GetStats(c *gin.Context) {
	idStr := c.Param("id_item_instance")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id_item_instance"})
		return
	}

	stats, err := h.service.GetStats(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrKerusakanItemInstanceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "item instance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
