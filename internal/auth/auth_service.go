package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user_api/internal/dto/auth"
)

type AuthService struct {
	repository *UserRepository
	jwtService *JWTService
}

func NewAuthService(repository *UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{repository: repository, jwtService: jwtService}
}

func (s *AuthService) Register(dto auth.RegisterDto) (string, error) {
	if dto.Email == "" || dto.Password == "" {
		return "", errors.New("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: string(hashedPassword),
	}

	err = s.repository.Create(user)
	if err != nil {
		return "", err
	}

	token, err := s.jwtService.GenerateToken(user, 1)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(dto auth.LoginDto) (string, error) {
	user, err := s.repository.GetByEmail(dto.Email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := s.jwtService.GenerateToken(user, 1)
	if err != nil {
		return "", err
	}

	return token, nil
}
