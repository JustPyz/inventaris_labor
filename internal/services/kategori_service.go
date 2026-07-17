package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateKategoriInput struct {
	Kategori string `json:"Kategori"`
}

type UpdateKategoriInput struct {
	Kategori string `json:"Kategori"`
}

type KategoriService struct {
	repo *repositories.KategoriRepository
}

func NewKategoriService(repo *repositories.KategoriRepository) *KategoriService {
	return &KategoriService{repo: repo}
}

func (s *KategoriService) List() ([]models.Kategori, error) {
	items := make([]models.Kategori, 0)
	if err := s.repo.List(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *KategoriService) Create(input CreateKategoriInput) (*models.Kategori, error) {
	kategori := strings.TrimSpace(input.Kategori)
	if kategori == "" {
		return nil, fmt.Errorf("kategori is required")
	}

	data := &models.Kategori{Kategori: kategori}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KategoriService) Get(id uint) (*models.Kategori, error) {
	data := &models.Kategori{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KategoriService) Update(id uint, input UpdateKategoriInput) (*models.Kategori, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	kategori := strings.TrimSpace(input.Kategori)
	if kategori == "" {
		return nil, fmt.Errorf("kategori is required")
	}

	data.Kategori = kategori
	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KategoriService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
