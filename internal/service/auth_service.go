package service

import (
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Login(req request.LoginRequest) (*response.AuthResponse, error)
	RefreshToken(refreshToken string) (*response.AuthResponse, error)
	ValidateToken(tokenString string) (string, error)
}

type authService struct {
	userRepo           repository.UserRepository
	jwtSecret          string
	refreshSecret      string
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret, refreshSecret string, accessExpire, refreshExpire time.Duration) AuthService {
	return &authService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		refreshSecret:      refreshSecret,
		accessTokenExpire:  accessExpire,
		refreshTokenExpire: refreshExpire,
	}
}

func (s *authService) Login(req request.LoginRequest) (*response.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %v", err)
	}

	// Check if user exists
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Convert user to response
	userResponse := response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    time.Now().Add(s.accessTokenExpire).Unix(), // Gunakan config
		User:         userResponse,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*response.AuthResponse, error) {
	// Validate refresh token
	userID, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Find user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	newAccessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	userResponse := response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &response.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    time.Now().Add(24 * time.Hour).Unix(),
		User:         userResponse,
	}, nil
}

func (s *authService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.accessTokenExpire).Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) generateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.refreshTokenExpire).Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

func (s *authService) validateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return userID, nil
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "access" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return userID, nil
}
