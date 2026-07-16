package handlers

import (
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JurusanHandler struct {
	service *services.JurusanService
}

func NewJurusanHandler(service *services.JurusanService) *JurusanHandler {
	return &JurusanHandler{service: service}
}

func (h *JurusanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch jurusan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *JurusanHandler) Create(c *gin.Context) {
	var input services.CreateJurusanInput
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

func (h *JurusanHandler) Get(c *gin.Context) {
	id, ok := parseJurusanID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "jurusan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch jurusan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *JurusanHandler) Update(c *gin.Context) {
	id, ok := parseJurusanID(c)
	if !ok {
		return
	}

	var input services.UpdateJurusanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "jurusan not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *JurusanHandler) Delete(c *gin.Context) {
	id, ok := parseJurusanID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "jurusan not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete jurusan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jurusan deleted"})
}

func parseJurusanID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
