package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateLaborInput struct {
	JurusanID uint   `json:"id_jurusan"`
	Labor     string `json:"labor"`
}

type UpdateLaborInput struct {
	JurusanID *uint   `json:"id_jurusan"`
	Labor     *string `json:"labor"`
}

type LaborService struct {
	repo        *repositories.LaborRepository
	jurusanRepo *repositories.JurusanRepository
}

func NewLaborService(repo *repositories.LaborRepository, jurusanRepo *repositories.JurusanRepository) *LaborService {
	return &LaborService{repo: repo, jurusanRepo: jurusanRepo}
}

func (s *LaborService) List() ([]models.Labor, error) {
	labor := make([]models.Labor, 0)
	if err := s.repo.List(&labor); err != nil {
		return nil, err
	}

	return labor, nil
}

func (s *LaborService) Create(input CreateLaborInput) (*models.Labor, error) {
	if input.JurusanID == 0 {
		return nil, fmt.Errorf("id_jurusan is required")
	}

	var jurusan models.Jurusan
	if err := s.jurusanRepo.FindByID(input.JurusanID, &jurusan); err != nil {
		return nil, fmt.Errorf("jurusan not found")
	}

	labor := strings.TrimSpace(input.Labor)
	if labor == "" {
		return nil, fmt.Errorf("labor is required")
	}

	data := &models.Labor{JurusanID: input.JurusanID, Labor: labor}
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

	if input.JurusanID != nil {
		if *input.JurusanID == 0 {
			return nil, fmt.Errorf("id_jurusan is required")
		}

		var jurusan models.Jurusan
		if err := s.jurusanRepo.FindByID(*input.JurusanID, &jurusan); err != nil {
			return nil, fmt.Errorf("jurusan not found")
		}
		data.JurusanID = *input.JurusanID
	}

	if input.Labor != nil {
		labor := strings.TrimSpace(*input.Labor)
		if labor == "" {
			return nil, fmt.Errorf("labor is required")
		}
		data.Labor = labor
	}

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
