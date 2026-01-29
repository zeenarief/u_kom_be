package service

import (
	"errors"
	"fmt"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"
)

type UserService interface {
	CreateUser(req request.UserCreateRequest) (*response.UserWithRoleResponse, error)
	GetUserByID(id string) (*response.UserWithRolesResponseAndPermissions, error)
	GetAllUsers() ([]response.UserWithRoleResponse, error)
	UpdateUser(id string, req request.UserUpdateRequest, currentUserID string, currentUserPermissions []string) (*response.UserWithRoleResponse, error) // Tambahkan parameter
	DeleteUser(id string, currentUserID string, currentUserPermissions []string) error                                                                  // Tambahkan parameter
	ChangePassword(id string, currentPassword, newPassword string, currentUserID string, currentUserPermissions []string) error                         // Tambahkan parameter
	GetProfile(userID string) (*response.ProfileResponse, error)
	SyncUserRoles(userID string, roleNames []string, currentUserID string, currentUserPermissions []string) error             // Tambahkan parameter
	SyncUserPermissions(userID string, permissionNames []string, currentUserID string, currentUserPermissions []string) error // Tambahkan parameter
	GetUserWithRolesAndPermissions(userID string) (*response.UserWithRolesResponseAndPermissions, error)
	GetUserPermissions(userID string) ([]string, error) // Tambahkan method baru
}

type userService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) UserService {
	return &userService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *userService) CreateUser(req request.UserCreateRequest) (*response.UserWithRoleResponse, error) {
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

	// Validasi Password Kuat
	if err := utils.ValidatePasswordComplexity(req.Password); err != nil {
		return nil, err // Akan melempar error ke handler
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Convert request to domain model
	user := &domain.User{
		ID:       utils.GenerateUUID(),
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

	// Handle role assignment
	if len(req.RoleIDs) > 0 {
		// Assign specified roles
		err = s.userRepo.SyncRoles(user.ID, req.RoleIDs)
		if err != nil {
			return nil, err
		}
	} else {
		// Assign default user role
		defaultRole, err := s.userRepo.GetDefaultRole()
		if err != nil {
			return nil, err
		}

		if defaultRole == nil {
			return nil, errors.New("default role not found")
		}

		err = s.userRepo.AssignRole(user.ID, defaultRole.ID)
		if err != nil {
			return nil, err
		}
	}

	// Reload user dengan roles dan permissions
	createdUser, err := s.userRepo.GetUserWithRolesAndPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	return converter.ToUserWithRoleResponse(createdUser), nil
}

func (s *userService) GetUserByID(id string) (*response.UserWithRolesResponseAndPermissions, error) {
	user, err := s.userRepo.GetUserWithRolesAndPermissions(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return converter.ToUserWithRolesResponseAndPermissions(user), nil
}

func (s *userService) GetAllUsers() ([]response.UserWithRoleResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.UserWithRoleResponse
	for _, user := range users {
		responses = append(responses, *converter.ToUserWithRoleResponse(&user))
	}

	return responses, nil
}

func (s *userService) UpdateUser(id string, req request.UserUpdateRequest, currentUserID string, currentUserPermissions []string) (*response.UserWithRoleResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 1. Validasi Basic: User update diri sendiri atau punya akses update others
	if id != currentUserID && !s.hasPermission(currentUserPermissions, "users.update.others") {
		return nil, errors.New("unauthorized: you can only update your own profile")
	}

	// 2. Update Basic Fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	// === 3. Update Roles ===
	// Cek apakah request mengirim role_ids (array tidak nil)
	if req.RoleIDs != nil {
		// SECURITY CHECK: Hanya user dengan permission 'users.manage_roles' yang boleh ganti role
		if !s.hasPermission(currentUserPermissions, "users.manage_roles") {
			return nil, errors.New("unauthorized: insufficient permissions to change roles")
		}

		// Lakukan Sync Role (Menggunakan ID, bukan Name)
		// Kita gunakan s.userRepo.SyncRoles karena menerima []string ID
		if err := s.userRepo.SyncRoles(id, req.RoleIDs); err != nil {
			return nil, fmt.Errorf("failed to update roles: %v", err)
		}
	}

	// === 4. Update Password (Reset by Admin) ===
	if req.Password != "" {
		// Cek Permission: Apakah user boleh mengganti password orang lain?
		// Jika update diri sendiri -> Boleh (tapi biasanya lewat endpoint change-password biar aman butuh password lama)
		// Jika update orang lain -> Harus punya permission 'users.change_password.others'

		if id != currentUserID && !s.hasPermission(currentUserPermissions, "users.change_password.others") {
			return nil, errors.New("unauthorized: insufficient permissions to reset password")
		}

		// Validasi Password Kuat
		if err := utils.ValidatePasswordComplexity(req.Password); err != nil {
			return nil, err // Akan melempar error ke handler
		}

		// Hash password baru
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	// Simpan perubahan data user basic
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	// Reload user dengan roles dan permissions terbaru untuk response
	updatedUser, err := s.userRepo.GetUserWithRolesAndPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	return converter.ToUserWithRoleResponse(updatedUser), nil
}

func (s *userService) DeleteUser(id string, currentUserID string, currentUserPermissions []string) error {
	// Validasi: user tidak bisa menghapus dirinya sendiri dan harus memiliki permission users.delete
	if id == currentUserID {
		return errors.New("cannot delete your own account")
	}

	if !s.hasPermission(currentUserPermissions, "users.delete") {
		return errors.New("unauthorized: insufficient permissions")
	}

	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) ChangePassword(id string, currentPassword, newPassword string, currentUserID string, currentUserPermissions []string) error {
	// Validasi: user hanya bisa mengubah password sendiri kecuali memiliki permission users.change_password.others
	if id != currentUserID && !s.hasPermission(currentUserPermissions, "users.change_password.others") {
		return errors.New("unauthorized: you can only change your own password")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Jika mengubah password orang lain, skip current password verification
	if id == currentUserID {
		if !utils.CheckPasswordHash(currentPassword, user.Password) {
			return errors.New("current password is incorrect")
		}
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

func (s *userService) SyncUserRoles(userID string, roleNames []string, currentUserID string, currentUserPermissions []string) error {
	if !s.hasPermission(currentUserPermissions, "users.manage_roles") {
		return errors.New("unauthorized: insufficient permissions")
	}

	// Convert role names to IDs
	var roleIDs []string
	for _, roleName := range roleNames {
		role, err := s.roleRepo.FindByName(roleName)
		if err != nil {
			return fmt.Errorf("error finding role: %s - %v", roleName, err)
		}
		if role == nil {
			return fmt.Errorf("role not found: %s", roleName)
		}
		roleIDs = append(roleIDs, role.ID)
	}

	return s.userRepo.SyncRoles(userID, roleIDs)
}

func (s *userService) SyncUserPermissions(userID string, permissionNames []string, currentUserID string, currentUserPermissions []string) error {
	if !s.hasPermission(currentUserPermissions, "users.manage_permissions") {
		return errors.New("unauthorized: insufficient permissions")
	}

	// Convert permission names to IDs
	var permissionIDs []string
	for _, permName := range permissionNames {
		permission, err := s.permissionRepo.FindByName(permName)
		if err != nil {
			return fmt.Errorf("error finding permission: %s - %v", permName, err)
		}
		if permission == nil {
			return fmt.Errorf("permission not found: %s", permName)
		}
		permissionIDs = append(permissionIDs, permission.ID)
	}

	return s.userRepo.SyncPermissions(userID, permissionIDs)
}

func (s *userService) GetUserWithRolesAndPermissions(userID string) (*response.UserWithRolesResponseAndPermissions, error) {
	user, err := s.userRepo.GetUserWithRolesAndPermissions(userID)
	if err != nil {
		return nil, err
	}

	return converter.ToUserWithRolesResponseAndPermissions(user), nil
}

func (s *userService) GetUserPermissions(userID string) ([]string, error) {
	user, err := s.userRepo.GetUserWithRolesAndPermissions(userID)
	if err != nil {
		return nil, err
	}

	var permissions []string

	// Ambil permissions dari roles
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission.Name)
		}
	}

	// Ambil permissions langsung
	for _, permission := range user.Permissions {
		permissions = append(permissions, permission.Name)
	}

	// Remove duplicates
	return utils.RemoveDuplicates(permissions), nil
}

// Helper function untuk mengecek permission
func (s *userService) hasPermission(permissions []string, requiredPermission string) bool {
	for _, perm := range permissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}
