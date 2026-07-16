package services

import (
	"fmt"
	"strings"

	"invela-be/internal/models"
	"invela-be/internal/repositories"
)

type CreateRoleInput struct {
	Role string `json:"role"`
}

type UpdateRoleInput struct {
	Role string `json:"role"`
}

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) List() ([]models.Role, error) {
	roles := make([]models.Role, 0)
	if err := s.repo.List(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) Create(input CreateRoleInput) (*models.Role, error) {
	role := strings.TrimSpace(input.Role)
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}

	data := &models.Role{Role: role}
	if err := s.repo.Create(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *RoleService) Get(id uint) (*models.Role, error) {
	data := &models.Role{}
	if err := s.repo.FindByID(id, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *RoleService) Update(id uint, input UpdateRoleInput) (*models.Role, error) {
	data, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	role := strings.TrimSpace(input.Role)
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}

	data.Role = role
	if err := s.repo.Update(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *RoleService) Delete(id uint) error {
	data, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(data)
}
