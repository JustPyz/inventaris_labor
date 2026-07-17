package handlers

import (
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LaborHandler struct {
	service *services.LaborService
}

func NewLaborHandler(service *services.LaborService) *LaborHandler {
	return &LaborHandler{service: service}
}

func (h *LaborHandler) List(c *gin.Context) {
	labor, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch labor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": labor})
}

func (h *LaborHandler) Create(c *gin.Context) {
	var input services.CreateLaborInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	labor, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": labor})
}

func (h *LaborHandler) Get(c *gin.Context) {
	id, ok := parseLaborID(c)
	if !ok {
		return
	}

	labor, err := h.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "labor not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch labor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": labor})
}

func (h *LaborHandler) Update(c *gin.Context) {
	id, ok := parseLaborID(c)
	if !ok {
		return
	}

	var input services.UpdateLaborInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	labor, err := h.service.Update(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "labor not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": labor})
}

func (h *LaborHandler) Delete(c *gin.Context) {
	id, ok := parseLaborID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "labor not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete labor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "labor deleted"})
}

func parseLaborID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
