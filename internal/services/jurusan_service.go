package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateJurusanInput struct {
	NamaJurusan string `json:"nama_jurusan"`
}

type UpdateJurusanInput struct {
	NamaJurusan string `json:"nama_jurusan"`
}

type JurusanService struct {
	repo *repositories.JurusanRepository
}

func NewJurusanService(repo *repositories.JurusanRepository) *JurusanService {
	return &JurusanService{repo: repo}
}

func (s *JurusanService) List() ([]models.Jurusan, error) {
	items := make([]models.Jurusan, 0)
	if err := s.repo.List(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *JurusanService) Create(input CreateJurusanInput) (*models.Jurusan, error) {
	nama := strings.TrimSpace(input.NamaJurusan)
	if nama == "" {
		return nil, fmt.Errorf("nama_jurusan is required")
	}

	data := &models.Jurusan{NamaJurusan: nama}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *JurusanService) Get(id uint) (*models.Jurusan, error) {
	data := &models.Jurusan{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *JurusanService) Update(id uint, input UpdateJurusanInput) (*models.Jurusan, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	nama := strings.TrimSpace(input.NamaJurusan)
	if nama == "" {
		return nil, fmt.Errorf("nama_jurusan is required")
	}

	data.NamaJurusan = nama
	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *JurusanService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
