package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"invela-be/internal/models"
	"invela-be/internal/repositories"

	"gorm.io/gorm"
)

var (
	ErrPeminjamanInvalidInput           = errors.New("peminjaman input is invalid")
	ErrPeminjamanItemInstanceNotFound   = errors.New("item instance not found")
	ErrPeminjamanNotFound               = errors.New("peminjaman not found")
)

var allowedPeminjamanStatus = map[string]struct{}{
	"aktif":               {},
	"selesai":             {},
	"melewati batas waktu": {},
}

type CreatePeminjamanInput struct {
	ItemInstanceID uint   `json:"id_item_instance"`
	NamaPeminjam   string `json:"nama_peminjam"`
	NomorTelepon   string `json:"nomor_telepon"`
	TanggalPinjam  string `json:"tanggal_pinjam"`
	TanggalKembali string `json:"tanggal_kembali"`
	Status         string `json:"status"`
}

type UpdatePeminjamanInput struct {
	ItemInstanceID *uint   `json:"id_item_instance"`
	NamaPeminjam   *string `json:"nama_peminjam"`
	NomorTelepon   *string `json:"nomor_telepon"`
	TanggalPinjam  *string `json:"tanggal_pinjam"`
	TanggalKembali *string `json:"tanggal_kembali"`
	Status         *string `json:"status"`
}

type PeminjamanService struct {
	repo *repositories.PeminjamanRepository
}

func NewPeminjamanService(repo *repositories.PeminjamanRepository) *PeminjamanService {
	return &PeminjamanService{repo: repo}
}

const dateLayout = "2006-01-02"

func (s *PeminjamanService) List() ([]models.Peminjaman, error) {
	items := make([]models.Peminjaman, 0)
	if err := s.repo.List(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *PeminjamanService) Create(input CreatePeminjamanInput) (*models.Peminjaman, error) {
	if input.ItemInstanceID == 0 {
		return nil, fmt.Errorf("%w: id_item_instance is required", ErrPeminjamanInvalidInput)
	}

	itemInstanceExists, err := s.repo.ItemInstanceExists(input.ItemInstanceID)
	if err != nil {
		return nil, err
	}
	if !itemInstanceExists {
		return nil, fmt.Errorf("%w: id_item_instance %d", ErrPeminjamanItemInstanceNotFound, input.ItemInstanceID)
	}

	namaPeminjam := strings.TrimSpace(input.NamaPeminjam)
	if namaPeminjam == "" {
		return nil, fmt.Errorf("%w: nama_peminjam is required", ErrPeminjamanInvalidInput)
	}

	nomorTelepon := strings.TrimSpace(input.NomorTelepon)
	if nomorTelepon == "" {
		return nil, fmt.Errorf("%w: nomor_telepon is required", ErrPeminjamanInvalidInput)
	}

	if input.TanggalPinjam == "" {
		return nil, fmt.Errorf("%w: tanggal_pinjam is required", ErrPeminjamanInvalidInput)
	}
	tanggalPinjam, err := time.Parse(dateLayout, input.TanggalPinjam)
	if err != nil {
		return nil, fmt.Errorf("%w: tanggal_pinjam must be in format YYYY-MM-DD", ErrPeminjamanInvalidInput)
	}

	if input.TanggalKembali == "" {
		return nil, fmt.Errorf("%w: tanggal_kembali is required", ErrPeminjamanInvalidInput)
	}
	tanggalKembali, err := time.Parse(dateLayout, input.TanggalKembali)
	if err != nil {
		return nil, fmt.Errorf("%w: tanggal_kembali must be in format YYYY-MM-DD", ErrPeminjamanInvalidInput)
	}

	if !tanggalKembali.After(tanggalPinjam) {
		return nil, fmt.Errorf("%w: tanggal_kembali must be after tanggal_pinjam", ErrPeminjamanInvalidInput)
	}

	status := strings.ToLower(strings.TrimSpace(input.Status))
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrPeminjamanInvalidInput)
	}

	if _, ok := allowedPeminjamanStatus[status]; !ok {
		return nil, fmt.Errorf("%w: status must be one of aktif, selesai, melewati batas waktu", ErrPeminjamanInvalidInput)
	}

	data := &models.Peminjaman{
		ItemInstanceID: input.ItemInstanceID,
		NamaPeminjam:   namaPeminjam,
		NomorTelepon:   nomorTelepon,
		TanggalPinjam:  tanggalPinjam,
		TanggalKembali: tanggalKembali,
		Status:         status,
	}

	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PeminjamanService) Get(id uint) (*models.Peminjaman, error) {
	data := &models.Peminjaman{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPeminjamanNotFound
		}

		return nil, err
	}

	return data, nil
}

func (s *PeminjamanService) Update(id uint, input UpdatePeminjamanInput) (*models.Peminjaman, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if input.ItemInstanceID != nil {
		if *input.ItemInstanceID == 0 {
			return nil, fmt.Errorf("%w: id_item_instance is required", ErrPeminjamanInvalidInput)
		}

		exists, err := s.repo.ItemInstanceExists(*input.ItemInstanceID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w: id_item_instance %d", ErrPeminjamanItemInstanceNotFound, *input.ItemInstanceID)
		}

		data.ItemInstanceID = *input.ItemInstanceID
	}

	if input.NamaPeminjam != nil {
		namaPeminjam := strings.TrimSpace(*input.NamaPeminjam)
		if namaPeminjam == "" {
			return nil, fmt.Errorf("%w: nama_peminjam is required", ErrPeminjamanInvalidInput)
		}
		data.NamaPeminjam = namaPeminjam
	}

	if input.NomorTelepon != nil {
		nomorTelepon := strings.TrimSpace(*input.NomorTelepon)
		if nomorTelepon == "" {
			return nil, fmt.Errorf("%w: nomor_telepon is required", ErrPeminjamanInvalidInput)
		}
		data.NomorTelepon = nomorTelepon
	}

	if input.TanggalPinjam != nil {
		tanggalPinjam, err := time.Parse(dateLayout, *input.TanggalPinjam)
		if err != nil {
			return nil, fmt.Errorf("%w: tanggal_pinjam must be in format YYYY-MM-DD", ErrPeminjamanInvalidInput)
		}
		data.TanggalPinjam = tanggalPinjam
	}

	if input.TanggalKembali != nil {
		tanggalKembali, err := time.Parse(dateLayout, *input.TanggalKembali)
		if err != nil {
			return nil, fmt.Errorf("%w: tanggal_kembali must be in format YYYY-MM-DD", ErrPeminjamanInvalidInput)
		}
		data.TanggalKembali = tanggalKembali
	}

	if !data.TanggalKembali.After(data.TanggalPinjam) {
		return nil, fmt.Errorf("%w: tanggal_kembali must be after tanggal_pinjam", ErrPeminjamanInvalidInput)
	}

	if input.Status != nil {
		status := strings.ToLower(strings.TrimSpace(*input.Status))
		if status == "" {
			return nil, fmt.Errorf("%w: status is required", ErrPeminjamanInvalidInput)
		}

		if _, ok := allowedPeminjamanStatus[status]; !ok {
			return nil, fmt.Errorf("%w: status must be one of aktif, selesai, melewati batas waktu", ErrPeminjamanInvalidInput)
		}

		data.Status = status
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PeminjamanService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
