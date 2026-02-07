package service

import (
	"fmt"
	"time"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/google/uuid"
)

type AuthService interface {
	Register(req request.UserCreateRequest) (*response.UserWithRoleResponse, error)
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

func (s *authService) Register(req request.UserCreateRequest) (*response.UserWithRoleResponse, error) {
	// 1. Validasi Duplikat Email
	existingUser, err := s.userRepo.FindByEmailWithRelations(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %v", err)
	}
	if existingUser != nil {
		return nil, apperrors.NewConflictError("email already exists")
	}

	// 2. Validasi Duplikat Username
	existingUser, err = s.userRepo.FindByUsernameWithRelations(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error checking username: %v", err)
	}
	if existingUser != nil {
		return nil, apperrors.NewConflictError("username already exists")
	}

	// 3. Validasi Kekuatan Password
	// Kita gunakan utils yang sudah dibuat sebelumnya
	if err := utils.ValidatePasswordComplexity(req.Password); err != nil {
		return nil, err
	}

	// 4. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. Convert request to domain model
	user := &domain.User{
		// Gunakan utils.GenerateUUID() agar konsisten, atau kosongkan biar GORM Hook yang handle
		ID:       utils.GenerateUUID(),
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// 6. Save to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// 7. Assign Default Role (PENTING untuk Public Register)
	defaultRole, err := s.userRepo.GetDefaultRole()
	if err != nil {
		return nil, err
	}
	if defaultRole == nil {
		// Fallback error jika belum di-seed
		return nil, apperrors.NewInternalError("registration failed: default role not configured")
	}

	err = s.userRepo.AssignRole(user.ID, defaultRole.ID)
	if err != nil {
		return nil, err
	}

	// 8. Reload user untuk response
	createdUser, err := s.userRepo.GetUserWithRolesAndPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	return converter.ToUserWithRoleResponse(createdUser), nil
}

func (s *authService) Login(req request.LoginRequest) (*response.AuthResponse, error) {
	var user *domain.User
	var err error

	// Coba sebagai email dulu
	user, err = s.userRepo.FindByEmailWithRelations(req.Login)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid login or password")
	}

	// Jika tidak ditemukan sebagai email, coba sebagai username
	if user == nil {
		user, err = s.userRepo.FindByUsernameWithRelations(req.Login)
		if err != nil {
			return nil, apperrors.NewUnauthorizedError("invalid login or password")
		}
		if user == nil {
			return nil, apperrors.NewUnauthorizedError("invalid login or password")
		}
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, apperrors.NewUnauthorizedError("invalid login or password")
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
	if err := s.userRepo.UpdateTokenHash(user.ID, &tokenHash); err != nil {
		return nil, err
	}

	userResponse := *converter.ToUserWithRoleResponse(user)

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
	return s.userRepo.UpdateTokenHash(userID, nil)
}

func (s *authService) RefreshToken(refreshToken string) (*response.AuthResponse, error) {
	// Validate refresh token
	userID, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid refresh token")
	}

	// Find user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user not found")
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
	if err := s.userRepo.UpdateTokenHash(user.ID, &tokenHash); err != nil {
		return nil, err
	}

	userResponse := response.UserWithRoleResponse{
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
		return "", apperrors.NewUnauthorizedError("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "access" {
		return "", apperrors.NewUnauthorizedError("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", apperrors.NewUnauthorizedError("invalid user ID in token")
	}

	// ✅ Check if user has logged out (token hash is empty)
	currentTokenHash, err := s.userRepo.GetTokenHash(userID)
	if err != nil {
		return "", apperrors.NewNotFoundError("user not found")
	}

	// Jika token hash kosong, berarti user sudah logout
	if currentTokenHash == nil {
		return "", apperrors.NewUnauthorizedError("token revoked - user logged out")
	}

	// Hash the incoming token and compare with stored hash
	incomingTokenHash := utils.HashToken(tokenString)
	if incomingTokenHash != *currentTokenHash {
		return "", apperrors.NewUnauthorizedError("token revoked - new login detected")
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
		return "", apperrors.NewUnauthorizedError("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return "", apperrors.NewUnauthorizedError("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", apperrors.NewUnauthorizedError("invalid user ID in token")
	}

	return userID, nil
}

func (s *authService) convertToResponse(user *domain.User) *response.UserWithRoleResponse {
	return &response.UserWithRoleResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
