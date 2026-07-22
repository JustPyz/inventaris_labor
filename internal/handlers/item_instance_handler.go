package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemInstanceHandler struct {
	service *services.ItemInstanceService
}

func NewItemInstanceHandler(service *services.ItemInstanceService) *ItemInstanceHandler {
	return &ItemInstanceHandler{service: service}
}

func (h *ItemInstanceHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch item instances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ItemInstanceHandler) Create(c *gin.Context) {
	var input services.CreateItemInstanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, services.ErrItemInstanceInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrItemInstancePerangkatNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_perangkat not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create item instance"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *ItemInstanceHandler) Get(c *gin.Context) {
	id, ok := parseItemInstanceID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrItemInstanceNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "item instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch item instance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ItemInstanceHandler) Update(c *gin.Context) {
	id, ok := parseItemInstanceID(c)
	if !ok {
		return
	}

	var input services.UpdateItemInstanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrItemInstanceNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "item instance not found"})
			return
		}

		if errors.Is(err, services.ErrItemInstanceInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrItemInstancePerangkatNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_perangkat not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update item instance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *ItemInstanceHandler) Delete(c *gin.Context) {
	id, ok := parseItemInstanceID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrItemInstanceNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "item instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete item instance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item instance deleted"})
}

func parseItemInstanceID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
