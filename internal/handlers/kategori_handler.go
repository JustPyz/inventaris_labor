package handlers

import (
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KategoriHandler struct {
	service *services.KategoriService
}

func NewKategoriHandler(service *services.KategoriService) *KategoriHandler {
	return &KategoriHandler{service: service}
}

func (h *KategoriHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *KategoriHandler) Create(c *gin.Context) {
	var input services.CreateKategoriInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *KategoriHandler) Get(c *gin.Context) {
	id, ok := parseKategoriID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *KategoriHandler) Update(c *gin.Context) {
	id, ok := parseKategoriID(c)
	if !ok {
		return
	}

	var input services.UpdateKategoriInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *KategoriHandler) Delete(c *gin.Context) {
	id, ok := parseKategoriID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kategori not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kategori deleted"})
}

func parseKategoriID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
