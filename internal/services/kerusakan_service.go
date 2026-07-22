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
	ErrKerusakanInvalidInput         = errors.New("kerusakan input is invalid")
	ErrKerusakanItemInstanceNotFound = errors.New("item instance not found")
	ErrKerusakanNotFound             = errors.New("kerusakan not found")
)

var allowedKerusakanStatus = map[string]struct{}{
	"butuh tindakan":    {},
	"sedang diperbaiki": {},
	"selesai":           {},
}

type CreateKerusakanInput struct {
	ItemInstanceID uint   `json:"id_item_instance"`
	Deskripsi      string `json:"deskripsi"`
	Status         string `json:"status"` // Optional, default is 'butuh tindakan'
}

type UpdateKerusakanInput struct {
	ItemInstanceID *uint   `json:"id_item_instance"`
	Deskripsi      *string `json:"deskripsi"`
	Status         *string `json:"status"`
}

type KerusakanService struct {
	repo repositories.KerusakanRepository
}

func NewKerusakanService(repo repositories.KerusakanRepository) *KerusakanService {
	return &KerusakanService{repo: repo}
}

func (s *KerusakanService) List() ([]models.Kerusakan, error) {
	return s.repo.FindAll()
}

func (s *KerusakanService) Create(userID uint, input CreateKerusakanInput) (*models.Kerusakan, error) {
	if input.ItemInstanceID == 0 {
		return nil, fmt.Errorf("%w: id_item_instance is required", ErrKerusakanInvalidInput)
	}

	itemInstanceExists, err := s.repo.ItemInstanceExists(input.ItemInstanceID)
	if err != nil {
		return nil, err
	}
	if !itemInstanceExists {
		return nil, fmt.Errorf("%w: id_item_instance %d", ErrKerusakanItemInstanceNotFound, input.ItemInstanceID)
	}

	deskripsi := strings.TrimSpace(input.Deskripsi)
	if deskripsi == "" {
		return nil, fmt.Errorf("%w: deskripsi is required", ErrKerusakanInvalidInput)
	}

	status := "butuh tindakan"
	if input.Status != "" {
		status = strings.ToLower(strings.TrimSpace(input.Status))
		if _, ok := allowedKerusakanStatus[status]; !ok {
			return nil, fmt.Errorf("%w: status must be one of 'butuh tindakan', 'sedang diperbaiki', 'selesai'", ErrKerusakanInvalidInput)
		}
	}

	kerusakan := &models.Kerusakan{
		ItemInstanceID: input.ItemInstanceID,
		UserID:         userID,
		Deskripsi:      deskripsi,
		Status:         status,
	}

	if err := s.repo.Create(kerusakan); err != nil {
		return nil, err
	}

	return kerusakan, nil
}

func (s *KerusakanService) Get(id uint) (*models.Kerusakan, error) {
	kerusakan, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKerusakanNotFound
		}
		return nil, err
	}
	return kerusakan, nil
}

func (s *KerusakanService) Update(id uint, input UpdateKerusakanInput) (*models.Kerusakan, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if input.ItemInstanceID != nil {
		if *input.ItemInstanceID == 0 {
			return nil, fmt.Errorf("%w: id_item_instance is required", ErrKerusakanInvalidInput)
		}

		exists, err := s.repo.ItemInstanceExists(*input.ItemInstanceID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: id_item_instance %d", ErrKerusakanItemInstanceNotFound, *input.ItemInstanceID)
		}

		data.ItemInstanceID = *input.ItemInstanceID
	}

	if input.Deskripsi != nil {
		deskripsi := strings.TrimSpace(*input.Deskripsi)
		if deskripsi == "" {
			return nil, fmt.Errorf("%w: deskripsi is required", ErrKerusakanInvalidInput)
		}
		data.Deskripsi = deskripsi
	}

	if input.Status != nil {
		status := strings.ToLower(strings.TrimSpace(*input.Status))
		if status == "" {
			return nil, fmt.Errorf("%w: status is required", ErrKerusakanInvalidInput)
		}

		if _, ok := allowedKerusakanStatus[status]; !ok {
			return nil, fmt.Errorf("%w: status must be one of 'butuh tindakan', 'sedang diperbaiki', 'selesai'", ErrKerusakanInvalidInput)
		}

		data.Status = status
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KerusakanService) Delete(id uint) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
