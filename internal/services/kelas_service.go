package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateKelasInput struct {
	JurusanID uint   `json:"id_jurusan"`
	Kelas     string `json:"kelas"`
}

type UpdateKelasInput struct {
	JurusanID *uint   `json:"id_jurusan"`
	Kelas     *string `json:"kelas"`
}

type KelasService struct {
	repo        *repositories.KelasRepository
	jurusanRepo *repositories.JurusanRepository
}

func NewKelasService(repo *repositories.KelasRepository, jurusanRepo *repositories.JurusanRepository) *KelasService {
	return &KelasService{repo: repo, jurusanRepo: jurusanRepo}
}

func (s *KelasService) List() ([]models.Kelas, error) {
	kelas := make([]models.Kelas, 0)
	if err := s.repo.List(&kelas); err != nil {
		return nil, err
	}

	return kelas, nil
}

func (s *KelasService) Create(input CreateKelasInput) (*models.Kelas, error) {
	if input.JurusanID == 0 {
		return nil, fmt.Errorf("id_jurusan is required")
	}

	var jurusan models.Jurusan
	if err := s.jurusanRepo.FindByID(input.JurusanID, &jurusan); err != nil {
		return nil, fmt.Errorf("jurusan not found")
	}

	kelas := strings.TrimSpace(input.Kelas)
	if kelas == "" {
		return nil, fmt.Errorf("kelas is required")
	}

	data := &models.Kelas{JurusanID: input.JurusanID, Kelas: kelas}
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

	if input.Kelas != nil {
		kelas := strings.TrimSpace(*input.Kelas)
		if kelas == "" {
			return nil, fmt.Errorf("kelas is required")
		}
		data.Kelas = kelas
	}

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
