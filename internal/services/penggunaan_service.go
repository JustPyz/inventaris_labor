package services

import (
	"errors"
	"fmt"

	"invela-be/internal/models"
	"invela-be/internal/repositories"

	"gorm.io/gorm"
)

var (
	ErrPenggunaanInvalidInput = errors.New("penggunaan input is invalid")
	ErrPenggunaanUserNotFound = errors.New("user not found")
	ErrPenggunaanLaborNotFound = errors.New("labor not found")
	ErrPenggunaanKelasNotFound = errors.New("kelas not found")
	ErrPenggunaanNotFound      = errors.New("penggunaan not found")
)

const (
	minJamPelajaran uint = 1
	maxJamPelajaran uint = 12
)

type CreatePenggunaanInput struct {
	UserID              uint `json:"id_user"`
	LaborID             uint `json:"id_labor"`
	KelasID             uint `json:"id_kelas"`
	JamPelajaranMulai   uint `json:"jam_pelajaran_mulai"`
	JamPelajaranSelesai uint `json:"jam_pelajaran_selesai"`
}

type PenggunaanService struct {
	repo *repositories.PenggunaanRepository
}

func NewPenggunaanService(repo *repositories.PenggunaanRepository) *PenggunaanService {
	return &PenggunaanService{repo: repo}
}

func (s *PenggunaanService) List() ([]models.Penggunaan, error) {
	items := make([]models.Penggunaan, 0)
	if err := s.repo.List(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *PenggunaanService) Create(input CreatePenggunaanInput) (*models.Penggunaan, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("%w: id_user is required", ErrPenggunaanInvalidInput)
	}

	userExists, err := s.repo.UserExists(input.UserID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("%w: id_user %d", ErrPenggunaanUserNotFound, input.UserID)
	}

	var user models.User
	if err := s.repo.FindUserByID(input.UserID, &user); err != nil {
		return nil, err
	}

	if input.LaborID == 0 {
		return nil, fmt.Errorf("%w: id_labor is required", ErrPenggunaanInvalidInput)
	}

	laborExists, err := s.repo.LaborExists(input.LaborID)
	if err != nil {
		return nil, err
	}
	if !laborExists {
		return nil, fmt.Errorf("%w: id_labor %d", ErrPenggunaanLaborNotFound, input.LaborID)
	}

	if input.KelasID == 0 {
		return nil, fmt.Errorf("%w: id_kelas is required", ErrPenggunaanInvalidInput)
	}

	kelasExists, err := s.repo.KelasExists(input.KelasID)
	if err != nil {
		return nil, err
	}
	if !kelasExists {
		return nil, fmt.Errorf("%w: id_kelas %d", ErrPenggunaanKelasNotFound, input.KelasID)
	}

	if input.JamPelajaranMulai < minJamPelajaran || input.JamPelajaranMulai > maxJamPelajaran {
		return nil, fmt.Errorf("%w: jam_pelajaran_mulai must be between %d and %d", ErrPenggunaanInvalidInput, minJamPelajaran, maxJamPelajaran)
	}

	if input.JamPelajaranSelesai < minJamPelajaran || input.JamPelajaranSelesai > maxJamPelajaran {
		return nil, fmt.Errorf("%w: jam_pelajaran_selesai must be between %d and %d", ErrPenggunaanInvalidInput, minJamPelajaran, maxJamPelajaran)
	}

	if input.JamPelajaranSelesai < input.JamPelajaranMulai {
		return nil, fmt.Errorf("%w: jam_pelajaran_selesai must be greater than or equal to jam_pelajaran_mulai", ErrPenggunaanInvalidInput)
	}

	data := &models.Penggunaan{
		UserID:              input.UserID,
		NamaPengguna:        user.Username,
		LaborID:             input.LaborID,
		KelasID:             input.KelasID,
		JamPelajaranMulai:   input.JamPelajaranMulai,
		JamPelajaranSelesai: input.JamPelajaranSelesai,
	}

	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	if err := s.repo.FindByID(data.ID, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PenggunaanService) Get(id uint) (*models.Penggunaan, error) {
	data := &models.Penggunaan{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPenggunaanNotFound
		}

		return nil, err
	}

	return data, nil
}

func (s *PenggunaanService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
