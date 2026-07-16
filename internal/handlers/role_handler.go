package handlers

import (
	"net/http"
	"strconv"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var input services.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	role, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

func (h *RoleHandler) Get(c *gin.Context) {
	id, ok := parseRoleID(c)
	if !ok {
		return
	}

	role, err := h.service.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "role not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, ok := parseRoleID(c)
	if !ok {
		return
	}

	var input services.UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	role, err := h.service.Update(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "role not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, ok := parseRoleID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "role not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role deleted"})
}

func parseRoleID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return 0, false
	}

	return uint(value), true
}
