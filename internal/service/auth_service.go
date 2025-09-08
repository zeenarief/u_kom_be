package service

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/utils"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(req request.UserCreateRequest) (*response.UserResponse, error)
	Login(req request.LoginRequest) (*response.AuthResponse, error)
	RefreshToken(refreshToken string) (*response.AuthResponse, error)
	ValidateToken(tokenString string) (string, error)
	Logout(userID string) error
	GetUserWithPermissions(userID string) (*domain.User, error)
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

func (s *authService) Register(req request.UserCreateRequest) (*response.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %v", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error checking username: %v", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Convert request to domain model
	user := &domain.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Assign default user role
	defaultRole, err := s.userRepo.GetDefaultRole()
	if err != nil {
		return nil, err
	}

	if defaultRole == nil {
		return nil, errors.New("default role not found")
	}

	// Assign default user role
	err = s.userRepo.AssignRole(user.ID, defaultRole.ID)
	if err != nil {
		return nil, err
	}

	// Convert domain model to response
	return s.convertToResponse(user), nil
}

func (s *authService) Login(req request.LoginRequest) (*response.AuthResponse, error) {
	var user *domain.User
	var err error

	// Coba sebagai email dulu
	user, err = s.userRepo.FindByEmail(req.Login)
	if err != nil {
		return nil, errors.New("invalid login or password")
	}

	// Jika tidak ditemukan sebagai email, coba sebagai username
	if user == nil {
		user, err = s.userRepo.FindByUsername(req.Login)
		if err != nil {
			return nil, errors.New("invalid login or password")
		}
		if user == nil {
			return nil, errors.New("invalid login or password")
		}
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid login or password")
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

	// Hash the access token and save to database
	tokenHash := utils.HashToken(accessToken)
	if err := s.userRepo.UpdateTokenHash(user.ID, tokenHash); err != nil {
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
		ExpiresIn:    time.Now().Add(s.accessTokenExpire).Unix(),
		User:         userResponse,
	}, nil
}

func (s *authService) Logout(userID string) error {
	// Set token hash menjadi empty string, sehingga token sekarang tidak valid
	return s.userRepo.UpdateTokenHash(userID, "")
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

	// ✅ Update token hash for the new token
	tokenHash := utils.HashToken(newAccessToken)
	if err := s.userRepo.UpdateTokenHash(user.ID, tokenHash); err != nil {
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
		ExpiresIn:    time.Now().Add(s.accessTokenExpire).Unix(),
		User:         userResponse,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	// Validate token signature first
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

	// ✅ Check if user has logged out (token hash is empty)
	currentTokenHash, err := s.userRepo.GetTokenHash(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Jika token hash kosong, berarti user sudah logout
	if currentTokenHash == "" {
		return "", errors.New("token revoked - user logged out")
	}

	// Hash the incoming token and compare with stored hash
	incomingTokenHash := utils.HashToken(tokenString)
	if incomingTokenHash != currentTokenHash {
		return "", errors.New("token revoked - new login detected")
	}

	return userID, nil
}

func (s *authService) GetUserWithPermissions(userID string) (*domain.User, error) {
	return s.userRepo.GetUserWithRolesAndPermissions(userID)
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

func (s *authService) convertToResponse(user *domain.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
