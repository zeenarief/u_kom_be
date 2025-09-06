package service

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(req request.UserCreateRequest) (*response.UserResponse, error)
	GetUserByID(id string) (*response.UserResponse, error)
	GetAllUsers() ([]response.UserResponse, error)
	UpdateUser(id string, req request.UserUpdateRequest) (*response.UserResponse, error)
	DeleteUser(id string) error
	ChangePassword(id string, currentPassword, newPassword string) error
	GetProfile(userID string) (*response.ProfileResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(req request.UserCreateRequest) (*response.UserResponse, error) {
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

	// Convert domain model to response
	return s.convertToResponse(user), nil
}

func (s *userService) GetUserByID(id string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.convertToResponse(user), nil
}

func (s *userService) GetAllUsers() ([]response.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.UserResponse
	for _, user := range users {
		responses = append(responses, *s.convertToResponse(&user))
	}

	return responses, nil
}

func (s *userService) UpdateUser(id string, req request.UserUpdateRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		// Check if new email already exists
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(user), nil
}

func (s *userService) DeleteUser(id string) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) ChangePassword(id string, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *userService) GetProfile(userID string) (*response.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Logika khusus profile bisa ditambahkan di sini
	profileComplete := user.Name != "" && user.Email != ""

	return &response.ProfileResponse{
		ID:              user.ID,
		Username:        user.Username,
		Name:            user.Name,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		ProfileComplete: profileComplete,
		AvatarURL:       "", // Anda bisa menambahkan logic untuk avatar
	}, nil
}

func (s *userService) convertToResponse(user *domain.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
