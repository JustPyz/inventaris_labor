package services

import (
	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type RiwayatPerbaikanService struct {
	repo repositories.RiwayatPerbaikanRepository
}

func NewRiwayatPerbaikanService(repo repositories.RiwayatPerbaikanRepository) *RiwayatPerbaikanService {
	return &RiwayatPerbaikanService{repo: repo}
}

// List mengembalikan semua riwayat perbaikan, diurut dari terbaru.
func (s *RiwayatPerbaikanService) List() ([]models.RiwayatPerbaikan, error) {
	return s.repo.FindAll()
}

// GetByKodeAsset mengembalikan riwayat perbaikan untuk kode asset tertentu.
func (s *RiwayatPerbaikanService) GetByKodeAsset(kodeAsset string) ([]models.RiwayatPerbaikan, error) {
	return s.repo.FindByKodeAsset(kodeAsset)
}
