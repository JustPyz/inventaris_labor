package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PeminjamanHandler struct {
	service *services.PeminjamanService
}

func NewPeminjamanHandler(service *services.PeminjamanService) *PeminjamanHandler {
	return &PeminjamanHandler{service: service}
}

func (h *PeminjamanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *PeminjamanHandler) Create(c *gin.Context) {
	var input services.CreatePeminjamanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, services.ErrPeminjamanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPeminjamanItemInstanceNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_item_instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create peminjaman"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *PeminjamanHandler) Get(c *gin.Context) {
	id, ok := parsePeminjamanID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrPeminjamanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "peminjaman not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PeminjamanHandler) Update(c *gin.Context) {
	id, ok := parsePeminjamanID(c)
	if !ok {
		return
	}

	var input services.UpdatePeminjamanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrPeminjamanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "peminjaman not found"})
			return
		}

		if errors.Is(err, services.ErrPeminjamanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPeminjamanItemInstanceNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_item_instance not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PeminjamanHandler) Delete(c *gin.Context) {
	id, ok := parsePeminjamanID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrPeminjamanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "peminjaman not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete peminjaman"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "peminjaman deleted"})
}

func parsePeminjamanID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
