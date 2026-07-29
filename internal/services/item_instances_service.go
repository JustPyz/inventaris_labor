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
	ErrItemInstanceInvalidInput      = errors.New("item instance input is invalid")
	ErrItemInstancePerangkatNotFound = errors.New("perangkat not found")
	ErrItemInstanceNotFound          = errors.New("item instance not found")
)

var allowedItemInstanceStatus = map[string]struct{}{
	"aktif":     {},
	"rusak":     {},
	"perbaikan": {},
	"nonaktif":  {},
}

type CreateItemInstanceInput struct {
	PerangkatID uint   `json:"id_perangkat"`
	KodeAsset   string `json:"kode_asset"`
	Status      string `json:"status"`
}

type UpdateItemInstanceInput struct {
	PerangkatID *uint   `json:"id_perangkat"`
	KodeAsset   *string `json:"kode_asset"`
	Status      *string `json:"status"`
}

type ItemInstanceService struct {
	repo *repositories.ItemInstanceRepository
}

func NewItemInstanceService(repo *repositories.ItemInstanceRepository) *ItemInstanceService {
	return &ItemInstanceService{repo: repo}
}

func (s *ItemInstanceService) List(role string, jurusanID *uint) ([]models.ItemInstance, error) {
	items := make([]models.ItemInstance, 0)
	if err := s.repo.List(role, jurusanID, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemInstanceService) Create(input CreateItemInstanceInput) (*models.ItemInstance, error) {
	if input.PerangkatID == 0 {
		return nil, fmt.Errorf("%w: id_perangkat is required", ErrItemInstanceInvalidInput)
	}

	perangkatExists, err := s.repo.PerangkatExists(input.PerangkatID)
	if err != nil {
		return nil, err
	}
	if !perangkatExists {
		return nil, fmt.Errorf("%w: id_perangkat %d", ErrItemInstancePerangkatNotFound, input.PerangkatID)
	}

	kodeAsset := strings.TrimSpace(input.KodeAsset)
	if kodeAsset == "" {
		return nil, fmt.Errorf("%w: kode_asset is required", ErrItemInstanceInvalidInput)
	}

	kodeAssetExists, err := s.repo.KodeAssetExists(kodeAsset)
	if err != nil {
		return nil, err
	}
	if kodeAssetExists {
		return nil, fmt.Errorf("%w: kode_asset already exists", ErrItemInstanceInvalidInput)
	}

	status := strings.ToLower(strings.TrimSpace(input.Status))
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrItemInstanceInvalidInput)
	}

	if _, ok := allowedItemInstanceStatus[status]; !ok {
		return nil, fmt.Errorf("%w: status must be one of aktif, rusak, perbaikan, nonaktif", ErrItemInstanceInvalidInput)
	}

	data := &models.ItemInstance{
		PerangkatID: input.PerangkatID,
		KodeAsset:   kodeAsset,
		Status:      status,
	}

	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ItemInstanceService) Get(id uint) (*models.ItemInstance, error) {
	data := &models.ItemInstance{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemInstanceNotFound
		}

		return nil, err
	}

	return data, nil
}

func (s *ItemInstanceService) Update(id uint, input UpdateItemInstanceInput) (*models.ItemInstance, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if input.PerangkatID != nil {
		if *input.PerangkatID == 0 {
			return nil, fmt.Errorf("%w: id_perangkat is required", ErrItemInstanceInvalidInput)
		}

		exists, err := s.repo.PerangkatExists(*input.PerangkatID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: id_perangkat %d", ErrItemInstancePerangkatNotFound, *input.PerangkatID)
		}

		data.PerangkatID = *input.PerangkatID
	}

	if input.KodeAsset != nil {
		kodeAsset := strings.TrimSpace(*input.KodeAsset)
		if kodeAsset == "" {
			return nil, fmt.Errorf("%w: kode_asset is required", ErrItemInstanceInvalidInput)
		}

		exists, err := s.repo.KodeAssetExistsExceptID(id, kodeAsset)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("%w: kode_asset already exists", ErrItemInstanceInvalidInput)
		}

		data.KodeAsset = kodeAsset
	}

	if input.Status != nil {
		status := strings.ToLower(strings.TrimSpace(*input.Status))
		if status == "" {
			return nil, fmt.Errorf("%w: status is required", ErrItemInstanceInvalidInput)
		}

		if _, ok := allowedItemInstanceStatus[status]; !ok {
			return nil, fmt.Errorf("%w: status must be one of aktif, rusak, perbaikan, nonaktif", ErrItemInstanceInvalidInput)
		}

		data.Status = status
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ItemInstanceService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
