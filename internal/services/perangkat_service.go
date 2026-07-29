package services

import (
	"errors"
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"

	"gorm.io/gorm"
)

var (
	ErrPerangkatInvalidInput     = errors.New("perangkat input is invalid")
	ErrPerangkatKategoriNotFound = errors.New("kategori not found")
	ErrPerangkatJurusanNotFound  = errors.New("jurusan not found")
	ErrPerangkatLaborNotFound    = errors.New("labor not found")
	ErrPerangkatNotFound         = errors.New("perangkat not found")
)

type CreatePerangkatInput struct {
	NamaPerangkat string  `json:"nama_perangkat"`
	KategoriID    uint    `json:"kategori_id"`
	JurusanID     uint    `json:"id_jurusan"`
	LaborID       uint    `json:"id_labor"`
	Deskripsi     *string `json:"deskripsi"`
}

type UpdatePerangkatInput struct {
	NamaPerangkat *string `json:"nama_perangkat"`
	KategoriID    *uint   `json:"kategori_id"`
	JurusanID     *uint   `json:"id_jurusan"`
	LaborID       *uint   `json:"id_labor"`
	Deskripsi     *string `json:"deskripsi"`
}

type PerangkatService struct {
	repo *repositories.PerangkatRepository
}

func NewPerangkatService(repo *repositories.PerangkatRepository) *PerangkatService {
	return &PerangkatService{repo: repo}
}

func (s *PerangkatService) List(role string, jurusanID *uint) ([]models.Perangkat, error) {
	items := make([]models.Perangkat, 0)
	if err := s.repo.List(role, jurusanID, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *PerangkatService) Create(input CreatePerangkatInput) (*models.Perangkat, error) {
	namaPerangkat := strings.TrimSpace(input.NamaPerangkat)
	if namaPerangkat == "" {
		return nil, fmt.Errorf("%w: nama_perangkat is required", ErrPerangkatInvalidInput)
	}

	if input.KategoriID == 0 {
		return nil, fmt.Errorf("%w: kategori_id is required", ErrPerangkatInvalidInput)
	}

	if input.JurusanID == 0 {
		return nil, fmt.Errorf("%w: id_jurusan is required", ErrPerangkatInvalidInput)
	}

	if input.LaborID == 0 {
		return nil, fmt.Errorf("%w: id_labor is required", ErrPerangkatInvalidInput)
	}

	kategoriExists, err := s.repo.KategoriExists(input.KategoriID)
	if err != nil {
		return nil, err
	}
	if !kategoriExists {
		return nil, fmt.Errorf("%w: kategori_id %d", ErrPerangkatKategoriNotFound, input.KategoriID)
	}

	jurusanExists, err := s.repo.JurusanExists(input.JurusanID)
	if err != nil {
		return nil, err
	}
	if !jurusanExists {
		return nil, fmt.Errorf("%w: id_jurusan %d", ErrPerangkatJurusanNotFound, input.JurusanID)
	}

	laborExists, err := s.repo.LaborExists(input.LaborID)
	if err != nil {
		return nil, err
	}
	if !laborExists {
		return nil, fmt.Errorf("%w: id_labor %d", ErrPerangkatLaborNotFound, input.LaborID)
	}

	deskripsi := "tidak ada deskripsi"
	if input.Deskripsi != nil && strings.TrimSpace(*input.Deskripsi) != "" {
		deskripsi = strings.TrimSpace(*input.Deskripsi)
	}

	data := &models.Perangkat{
		NamaPerangkat: namaPerangkat,
		KategoriID:    input.KategoriID,
		JurusanID:     input.JurusanID,
		LaborID:       input.LaborID,
		Deskripsi:     deskripsi,
	}

	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PerangkatService) Get(id uint) (*models.Perangkat, error) {
	data := &models.Perangkat{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPerangkatNotFound
		}

		return nil, err
	}

	return data, nil
}

func (s *PerangkatService) Update(id uint, input UpdatePerangkatInput) (*models.Perangkat, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if input.NamaPerangkat != nil {
		namaPerangkat := strings.TrimSpace(*input.NamaPerangkat)
		if namaPerangkat == "" {
			return nil, fmt.Errorf("%w: nama_perangkat is required", ErrPerangkatInvalidInput)
		}
		data.NamaPerangkat = namaPerangkat
	}

	if input.KategoriID != nil {
		if *input.KategoriID == 0 {
			return nil, fmt.Errorf("%w: kategori_id is required", ErrPerangkatInvalidInput)
		}

		exists, err := s.repo.KategoriExists(*input.KategoriID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: kategori_id %d", ErrPerangkatKategoriNotFound, *input.KategoriID)
		}

		data.KategoriID = *input.KategoriID
	}

	if input.JurusanID != nil {
		if *input.JurusanID == 0 {
			return nil, fmt.Errorf("%w: id_jurusan is required", ErrPerangkatInvalidInput)
		}

		exists, err := s.repo.JurusanExists(*input.JurusanID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: id_jurusan %d", ErrPerangkatJurusanNotFound, *input.JurusanID)
		}

		data.JurusanID = *input.JurusanID
	}

	if input.LaborID != nil {
		if *input.LaborID == 0 {
			return nil, fmt.Errorf("%w: id_labor is required", ErrPerangkatInvalidInput)
		}

		exists, err := s.repo.LaborExists(*input.LaborID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: id_labor %d", ErrPerangkatLaborNotFound, *input.LaborID)
		}

		data.LaborID = *input.LaborID
	}

	if input.Deskripsi != nil {
		deskripsi := strings.TrimSpace(*input.Deskripsi)
		if deskripsi == "" {
			deskripsi = "tidak ada deskripsi"
		}
		data.Deskripsi = deskripsi
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PerangkatService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
