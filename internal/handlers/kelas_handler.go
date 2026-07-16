package handlers

import (
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KelasHandler struct {
	service *services.KelasService
}

func NewKelasHandler(service *services.KelasService) *KelasHandler {
	return &KelasHandler{service: service}
}

func (h *KelasHandler) List(c *gin.Context) {
	kelas, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kelas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kelas})
}

func (h *KelasHandler) Create(c *gin.Context) {
	var input services.CreateKelasInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	kelas, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": kelas})
}

func (h *KelasHandler) Get(c *gin.Context) {
	id, ok := parseKelasID(c)
	if !ok {
		return
	}

	kelas, err := h.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kelas not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch kelas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kelas})
}

func (h *KelasHandler) Update(c *gin.Context) {
	id, ok := parseKelasID(c)
	if !ok {
		return
	}

	var input services.UpdateKelasInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	kelas, err := h.service.Update(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kelas not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kelas})
}

func (h *KelasHandler) Delete(c *gin.Context) {
	id, ok := parseKelasID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "kelas not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete kelas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kelas deleted"})
}

func parseKelasID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
