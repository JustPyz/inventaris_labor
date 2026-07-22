package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenggunaanHandler struct {
	service *services.PenggunaanService
}

func NewPenggunaanHandler(service *services.PenggunaanService) *PenggunaanHandler {
	return &PenggunaanHandler{service: service}
}

func (h *PenggunaanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch penggunaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *PenggunaanHandler) Create(c *gin.Context) {
	var input services.CreatePenggunaanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, services.ErrPenggunaanInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPenggunaanUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_user not found"})
			return
		}

		if errors.Is(err, services.ErrPenggunaanLaborNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_labor not found"})
			return
		}

		if errors.Is(err, services.ErrPenggunaanKelasNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_kelas not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create penggunaan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *PenggunaanHandler) Get(c *gin.Context) {
	id, ok := parsePenggunaanID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrPenggunaanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "penggunaan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch penggunaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PenggunaanHandler) Delete(c *gin.Context) {
	id, ok := parsePenggunaanID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrPenggunaanNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "penggunaan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete penggunaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "penggunaan deleted"})
}

func parsePenggunaanID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
