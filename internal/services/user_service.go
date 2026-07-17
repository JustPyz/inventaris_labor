package services

import (
	"fmt"
	"strings"
	"time"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateUserInput struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	RoleID       uint   `json:"role_id"`
}

type UpdateUserInput struct {
	Username     *string `json:"username"`
	PasswordHash *string `json:"password_hash"`
	RoleID       *uint   `json:"role_id"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
		return nil, fmt.Errorf("username is required")
	}

	passwordHash := strings.TrimSpace(input.PasswordHash)
	if passwordHash == "" {
		return nil, fmt.Errorf("password_hash is required")
	}

	if input.RoleID == 0 {
		return nil, fmt.Errorf("role_id is required")
	}

	data := &models.User{
		Username:     username,
		PasswordHash: passwordHash,
		RoleID:       input.RoleID,
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
		return nil, err
	}

	response := mapUserToResponse(*data)
	return &response, nil
}

func (s *UserService) Update(id uint, input UpdateUserInput) (*UserResponse, error) {
	data := &models.User{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	if input.Username != nil {
		username := strings.TrimSpace(*input.Username)
		if username == "" {
			return nil, fmt.Errorf("username is required")
		}
		data.Username = username
	}

	if input.PasswordHash != nil {
		passwordHash := strings.TrimSpace(*input.PasswordHash)
		if passwordHash == "" {
			return nil, fmt.Errorf("password_hash is required")
		}
		data.PasswordHash = passwordHash
	}

	if input.RoleID != nil {
		if *input.RoleID == 0 {
			return nil, fmt.Errorf("role_id is required")
		}
		data.RoleID = *input.RoleID
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
	return UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
