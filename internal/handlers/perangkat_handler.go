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

type PerangkatHandler struct {
	service *services.PerangkatService
}

func NewPerangkatHandler(service *services.PerangkatService) *PerangkatHandler {
	return &PerangkatHandler{service: service}
}

func (h *PerangkatHandler) List(c *gin.Context) {
	roleStr := ""
	if role, exists := c.Get(middlewares.ContextKeyRole); exists {
		if r, ok := role.(string); ok {
			roleStr = r
		}
	}

	var jurusanID *uint
	if jID, exists := c.Get(middlewares.ContextKeyJurusanID); exists {
		if j, ok := jID.(*uint); ok {
			jurusanID = j
		}
	}

	items, err := h.service.List(roleStr, jurusanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch perangkat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *PerangkatHandler) Create(c *gin.Context) {
	var input services.CreatePerangkatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, services.ErrPerangkatInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPerangkatKategoriNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "kategori_id not found"})
			return
		}

		if errors.Is(err, services.ErrPerangkatJurusanNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_jurusan not found"})
			return
		}

		if errors.Is(err, services.ErrPerangkatLaborNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_labor not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create perangkat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *PerangkatHandler) Get(c *gin.Context) {
	id, ok := parsePerangkatID(c)
	if !ok {
		return
	}

	item, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, services.ErrPerangkatNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perangkat not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch perangkat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PerangkatHandler) Update(c *gin.Context) {
	id, ok := parsePerangkatID(c)
	if !ok {
		return
	}

	var input services.UpdatePerangkatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, services.ErrPerangkatNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perangkat not found"})
			return
		}

		if errors.Is(err, services.ErrPerangkatInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPerangkatKategoriNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "kategori_id not found"})
			return
		}

		if errors.Is(err, services.ErrPerangkatJurusanNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_jurusan not found"})
			return
		}

		if errors.Is(err, services.ErrPerangkatLaborNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "id_labor not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update perangkat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *PerangkatHandler) Delete(c *gin.Context) {
	id, ok := parsePerangkatID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, services.ErrPerangkatNotFound) || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "perangkat not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete perangkat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "perangkat deleted"})
}

func parsePerangkatID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
