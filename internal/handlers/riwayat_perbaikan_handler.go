package handlers

import (
	"net/http"
	"strings"

	"invela-be/internal/services"

	"github.com/gin-gonic/gin"
)

type RiwayatPerbaikanHandler struct {
	service *services.RiwayatPerbaikanService
}

func NewRiwayatPerbaikanHandler(service *services.RiwayatPerbaikanService) *RiwayatPerbaikanHandler {
	return &RiwayatPerbaikanHandler{service: service}
}

// List mengembalikan seluruh riwayat perbaikan, diurut dari terbaru.
func (h *RiwayatPerbaikanHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch riwayat perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetByKodeAsset mengembalikan riwayat perbaikan untuk satu kode asset.
func (h *RiwayatPerbaikanHandler) GetByKodeAsset(c *gin.Context) {
	kodeAsset := strings.TrimSpace(c.Param("kode_asset"))
	if kodeAsset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "kode_asset is required"})
		return
	}

	items, err := h.service.GetByKodeAsset(kodeAsset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch riwayat perbaikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}
