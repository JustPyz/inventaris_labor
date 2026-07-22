package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"invela-be/internal/models"
	"invela-be/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserInvalidInput = errors.New("user input is invalid")
	ErrRoleNotFound     = errors.New("role not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrJurusanNotFound  = errors.New("jurusan not found")
)

type CreateUserInput struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	RoleID       uint   `json:"role_id"`
	JurusanID    *uint  `json:"jurusan_id"`
}

type UpdateUserInput struct {
	Username     *string `json:"username"`
	PasswordHash *string `json:"password_hash"`
	RoleID       *uint   `json:"role_id"`
	JurusanID    *uint   `json:"jurusan_id"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	JurusanID *uint     `json:"jurusan_id"`
	Jurusan   *string   `json:"jurusan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List() ([]UserResponse, error) {
	items := make([]models.User, 0)
	if err := s.repo.List(&items); err != nil {
		return nil, err
	}

	return mapUsersToResponses(items), nil
}

func (s *UserService) Create(input CreateUserInput) (*UserResponse, error) {
	username := strings.TrimSpace(input.Username)
	if username == "" {
		return nil, fmt.Errorf("%w: username is required", ErrUserInvalidInput)
	}

	passwordRaw := strings.TrimSpace(input.PasswordHash)
	if passwordRaw == "" {
		return nil, fmt.Errorf("%w: password_hash is required", ErrUserInvalidInput)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if input.RoleID == 0 {
		return nil, fmt.Errorf("%w: role_id is required", ErrUserInvalidInput)
	}

	roleExists, err := s.repo.RoleExists(input.RoleID)
	if err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, fmt.Errorf("%w: role_id %d", ErrRoleNotFound, input.RoleID)
	}

	var jurusanID *uint
	if input.JurusanID != nil && *input.JurusanID != 0 {
		jurusanExists, err := s.repo.JurusanExists(*input.JurusanID)
		if err != nil {
			return nil, err
		}
		if !jurusanExists {
			return nil, fmt.Errorf("%w: jurusan_id %d", ErrJurusanNotFound, *input.JurusanID)
		}
		jurusanID = input.JurusanID
	}

	data := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		RoleID:       input.RoleID,
		JurusanID:    jurusanID,
	}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	if err := s.repo.FindByID(data.ID, data); err != nil {
		return nil, err
	}

	response := mapUserToResponse(*data)
	return &response, nil
}

func (s *UserService) Get(id uint) (*UserResponse, error) {
	data := &models.User{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	response := mapUserToResponse(*data)
	return &response, nil
}

func (s *UserService) Update(id uint, input UpdateUserInput) (*UserResponse, error) {
	data := &models.User{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if input.Username != nil {
		username := strings.TrimSpace(*input.Username)
		if username == "" {
			return nil, fmt.Errorf("%w: username is required", ErrUserInvalidInput)
		}
		data.Username = username
	}

	if input.PasswordHash != nil {
		passwordRaw := strings.TrimSpace(*input.PasswordHash)
		if passwordRaw == "" {
			return nil, fmt.Errorf("%w: password_hash is required", ErrUserInvalidInput)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		data.PasswordHash = string(hashedPassword)
	}

	if input.RoleID != nil {
		if *input.RoleID == 0 {
			return nil, fmt.Errorf("%w: role_id is required", ErrUserInvalidInput)
		}

		roleExists, err := s.repo.RoleExists(*input.RoleID)
		if err != nil {
			return nil, err
		}
		if !roleExists {
			return nil, fmt.Errorf("%w: role_id %d", ErrRoleNotFound, *input.RoleID)
		}
		data.RoleID = *input.RoleID
	}

	if input.JurusanID != nil {
		if *input.JurusanID == 0 {
			data.JurusanID = nil
		} else {
			jurusanExists, err := s.repo.JurusanExists(*input.JurusanID)
			if err != nil {
				return nil, err
			}
			if !jurusanExists {
				return nil, fmt.Errorf("%w: jurusan_id %d", ErrJurusanNotFound, *input.JurusanID)
			}
			data.JurusanID = input.JurusanID
		}
	}

	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	if err := s.repo.FindByID(data.ID, data); err != nil {
		return nil, err
	}

	response := mapUserToResponse(*data)
	return &response, nil
}

func (s *UserService) Delete(id uint) error {
	data, err := s.GetModel(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}

func (s *UserService) GetModel(id uint) (*models.User, error) {
	data := &models.User{}
	if err := s.repo.FindByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return data, nil
}

func mapUsersToResponses(users []models.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, mapUserToResponse(user))
	}

	return responses
}

func mapUserToResponse(user models.User) UserResponse {
	var jurusanNama *string
	if user.Jurusan != nil {
		jurusanNama = &user.Jurusan.NamaJurusan
	}

	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role.Role,
		JurusanID: user.JurusanID,
		Jurusan:   jurusanNama,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
