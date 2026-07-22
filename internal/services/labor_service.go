package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateLaborInput struct {
	Labor string `json:"labor"`
}

type UpdateLaborInput struct {
	Labor string `json:"labor"`
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
	labor := strings.TrimSpace(input.Labor)
	if labor == "" {
		return nil, fmt.Errorf("labor is required")
	}

	data := &models.Labor{Labor: labor}
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

	labor := strings.TrimSpace(input.Labor)
	if labor == "" {
		return nil, fmt.Errorf("labor is required")
	}

	data.Labor = labor
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
