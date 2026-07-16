package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateKelasInput struct {
	Kelas string `json:"kelas"`
}

type UpdateKelasInput struct {
	Kelas string `json:"kelas"`
}

type KelasService struct {
	repo *repositories.KelasRepository
}

func NewKelasService(repo *repositories.KelasRepository) *KelasService {
	return &KelasService{repo: repo}
}

func (s *KelasService) List() ([]models.Kelas, error) {
	kelas := make([]models.Kelas, 0)
	if err := s.repo.List(&kelas); err != nil {
		return nil, err
	}

	return kelas, nil
}

func (s *KelasService) Create(input CreateKelasInput) (*models.Kelas, error) {
	kelas := strings.TrimSpace(input.Kelas)
	if kelas == "" {
		return nil, fmt.Errorf("kelas is required")
	}

	data := &models.Kelas{Kelas: kelas}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KelasService) Get(id uint) (*models.Kelas, error) {
	data := &models.Kelas{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KelasService) Update(id uint, input UpdateKelasInput) (*models.Kelas, error) {
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

func (s *KelasService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
