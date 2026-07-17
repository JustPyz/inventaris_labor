package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateLaborInput struct {
	Kelas string `json:"kelas"`
}

type UpdateLaborInput struct {
	Kelas string `json:"kelas"`
}

type LaborService struct {
	repo *repositories.LaborRepository
}

func NewLaborService(repo *repositories.LaborRepository) *LaborService {
	return &LaborService{repo: repo}
}

func (s *LaborService) List() ([]models.Labor, error) {
	labor := make([]models.Labor, 0)
	if err := s.repo.List(&labor); err != nil {
		return nil, err
	}

	return labor, nil
}

func (s *LaborService) Create(input CreateLaborInput) (*models.Labor, error) {
	kelas := strings.TrimSpace(input.Kelas)
	if kelas == "" {
		return nil, fmt.Errorf("kelas is required")
	}

	data := &models.Labor{Kelas: kelas}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *LaborService) Get(id uint) (*models.Labor, error) {
	data := &models.Labor{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *LaborService) Update(id uint, input UpdateLaborInput) (*models.Labor, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	kelas := strings.TrimSpace(input.Kelas)
	if kelas == "" {
		return nil, fmt.Errorf("kelas is required")
	}

	data.Kelas = kelas
	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *LaborService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
