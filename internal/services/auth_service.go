package services

import (
	"errors"
	"strings"
	"time"

	"invela-be/internal/models"
	"invela-be/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token string     `json:"token"`
	User  UserInfo   `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Login(input LoginInput) (*LoginResult, error) {
	username := strings.TrimSpace(input.Username)
	password := strings.TrimSpace(input.Password)

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	var user models.User
	if err := s.userRepo.FindByUsername(username, &user); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role.Role,
		},
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
