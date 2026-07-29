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
	ErrPerbaikanInvalidInput      = errors.New("perbaikan input is invalid")
	ErrPerbaikanKerusakanNotFound = errors.New("kerusakan not found")
	ErrPerbaikanNotFound          = errors.New("perbaikan not found")
)

type CreatePerbaikanInput struct {
	KerusakanID        uint   `json:"id_kerusakan"`
	DeskripsiPerbaikan string `json:"deskripsi_perbaikan"`
	Biaya              int64  `json:"biaya"`
}

type UpdatePerbaikanInput struct {
	DeskripsiPerbaikan *string `json:"deskripsi_perbaikan"`
	Biaya              *int64  `json:"biaya"`
}

type PerbaikanService struct {
	db                    *gorm.DB
	repo                  repositories.PerbaikanRepository
	kerusakanRepo         repositories.KerusakanRepository
	itemInstanceRepo      *repositories.ItemInstanceRepository
	userRepo              *repositories.UserRepository
	riwayatPerbaikanRepo  repositories.RiwayatPerbaikanRepository
}

func NewPerbaikanService(
	db *gorm.DB,
	repo repositories.PerbaikanRepository,
	kerusakanRepo repositories.KerusakanRepository,
	itemInstanceRepo *repositories.ItemInstanceRepository,
	userRepo *repositories.UserRepository,
	riwayatPerbaikanRepo repositories.RiwayatPerbaikanRepository,
) *PerbaikanService {
	return &PerbaikanService{
		db:                   db,
		repo:                 repo,
		kerusakanRepo:        kerusakanRepo,
		itemInstanceRepo:     itemInstanceRepo,
		userRepo:             userRepo,
		riwayatPerbaikanRepo: riwayatPerbaikanRepo,
	}
}

func (s *PerbaikanService) List() ([]models.Perbaikan, error) {
	return s.repo.FindAll()
}

func (s *PerbaikanService) Create(userID uint, input CreatePerbaikanInput) (*models.Perbaikan, error) {
	if input.KerusakanID == 0 {
		return nil, fmt.Errorf("%w: id_kerusakan is required", ErrPerbaikanInvalidInput)
	}

	kerusakanExists, err := s.repo.KerusakanExists(input.KerusakanID)
	if err != nil {
		return nil, err
	}
	if !kerusakanExists {
		return nil, fmt.Errorf("%w: id_kerusakan %d", ErrPerbaikanKerusakanNotFound, input.KerusakanID)
	}

	deskripsi := strings.TrimSpace(input.DeskripsiPerbaikan)
	if deskripsi == "" {
		return nil, fmt.Errorf("%w: deskripsi_perbaikan is required", ErrPerbaikanInvalidInput)
	}

	if input.Biaya < 0 {
		return nil, fmt.Errorf("%w: biaya cannot be negative", ErrPerbaikanInvalidInput)
	}

	perbaikan := &models.Perbaikan{
		KerusakanID:        input.KerusakanID,
		UserID:             userID,
		DeskripsiPerbaikan: deskripsi,
		Biaya:              input.Biaya,
	}

	// Ambil data kerusakan sebelum transaksi untuk mendapatkan id_item_instance
	kerusakan, err := s.kerusakanRepo.FindByID(input.KerusakanID)
	if err != nil {
		return nil, fmt.Errorf("kerusakan tidak ditemukan: %w", err)
	}
	itemInstanceID := kerusakan.ItemInstanceID

	// Ambil item_instance beserta data Perangkat (sudah di-preload)
	var itemInstance models.ItemInstance
	if err := s.itemInstanceRepo.FindByID(itemInstanceID, &itemInstance); err != nil {
		return nil, fmt.Errorf("item instance tidak ditemukan: %w", err)
	}

	// Ambil data teknisi berdasarkan userID
	var user models.User
	if err := s.userRepo.FindByID(userID, &user); err != nil {
		return nil, fmt.Errorf("user tidak ditemukan: %w", err)
	}

	// Jalankan dalam transaksi: simpan perbaikan + update status kerusakan
	// + (kondisional) update status item_instance ke 'aktif'
	// + insert snapshot ke tabel riwayat_perbaikans.
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Langkah 1: Simpan data perbaikan
		if err := tx.Create(perbaikan).Error; err != nil {
			return fmt.Errorf("simpan perbaikan: %w", err)
		}

		// Langkah 2: Set status kerusakan → 'selesai'
		if err := tx.Model(&models.Kerusakan{}).Where("id = ?", input.KerusakanID).
			Update("status", "selesai").Error; err != nil {
			return fmt.Errorf("update status kerusakan: %w", err)
		}

		// Langkah 3: Cek apakah masih ada kerusakan lain yang belum 'selesai'
		// pada item_instance yang sama. Jika tidak ada, set item_instance → 'aktif'.
		hasOpen, err := s.kerusakanRepo.HasOpenKerusakan(itemInstanceID, input.KerusakanID)
		if err != nil {
			return fmt.Errorf("cek kerusakan terbuka: %w", err)
		}
		if !hasOpen {
			if err := tx.Model(&models.ItemInstance{}).Where("id = ?", itemInstanceID).
				Update("status", "aktif").Error; err != nil {
				return fmt.Errorf("update status item_instance: %w", err)
			}
		}

		// Langkah 4: Insert snapshot ke tabel riwayat_perbaikans.
		// Semua data adalah salinan nilai saat ini — bukan referensi FK.
		riwayat := &models.RiwayatPerbaikan{
			PerbaikanID:        perbaikan.ID,
			KerusakanID:        input.KerusakanID,
			KodeAsset:          itemInstance.KodeAsset,
			NamaPerangkat:      itemInstance.Perangkat.NamaPerangkat,
			DeskripsiKerusakan: kerusakan.Deskripsi,
			DeskripsiPerbaikan: deskripsi,
			Biaya:              input.Biaya,
			NamaTeknisi:        user.Username,
			TanggalPerbaikan:   perbaikan.CreatedAt,
		}
		if err := tx.Create(riwayat).Error; err != nil {
			return fmt.Errorf("simpan riwayat perbaikan: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return perbaikan, nil
}

func (s *PerbaikanService) Get(id uint) (*models.Perbaikan, error) {
	perbaikan, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPerbaikanNotFound
		}
		return nil, err
	}
	return perbaikan, nil
}

func (s *PerbaikanService) Update(id uint, input UpdatePerbaikanInput) (*models.Perbaikan, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if input.DeskripsiPerbaikan != nil {
		deskripsi := strings.TrimSpace(*input.DeskripsiPerbaikan)
		if deskripsi == "" {
			return nil, fmt.Errorf("%w: deskripsi_perbaikan is required", ErrPerbaikanInvalidInput)
		}
		data.DeskripsiPerbaikan = deskripsi
	}

	if input.Biaya != nil {
		if *input.Biaya < 0 {
			return nil, fmt.Errorf("%w: biaya cannot be negative", ErrPerbaikanInvalidInput)
		}
		data.Biaya = *input.Biaya
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PerbaikanService) Delete(id uint) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
