package handler

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Registration failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	CreatedResponse(c, "User created successfully", user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		NotFoundError(c, "User not found")
		return
	}

	SuccessResponse(c, "User retrieved successfully", user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Users retrieved successfully", users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req request.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Dapatkan current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		UnauthorizedError(c, "user not found in context")
		return
	}

	currentUserDomain := currentUser.(*domain.User)

	// Dapatkan permissions current user
	currentPermissions, err := h.userService.GetUserPermissions(currentUserDomain.ID)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user permissions"})
		InternalServerError(c, err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateUser(id, req, currentUserDomain.ID, currentPermissions)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		BadRequestError(c, "Bad request", err.Error())
		return
	}

	//c.JSON(http.StatusOK, updatedUser)
	SuccessResponse(c, "User updated successfully", updatedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Dapatkan current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		UnauthorizedError(c, "user not found in context")
		return
	}
	currentUserDomain := currentUser.(*domain.User)

	// Dapatkan permissions current user
	currentPermissions, err := h.userService.GetUserPermissions(currentUserDomain.ID)
	if err != nil {
		InternalServerError(c, "failed to get user permissions")
		return
	}

	err = h.userService.DeleteUser(id, currentUserDomain.ID, currentPermissions)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Dapatkan current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		UnauthorizedError(c, "user not found in context")
		return
	}
	currentUserDomain := currentUser.(*domain.User)

	// Dapatkan permissions current user
	currentPermissions, err := h.userService.GetUserPermissions(currentUserDomain.ID)
	if err != nil {
		InternalServerError(c, "failed to get user permissions")
		return
	}

	err = h.userService.ChangePassword(id, req.CurrentPassword, req.NewPassword, currentUserDomain.ID, currentPermissions)
	if err != nil {
		ForbiddenError(c, err.Error())
		return
	}

	SuccessResponse(c, "Password changed successfully", nil)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		UnauthorizedError(c, "User ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		InternalServerError(c, "Invalid user ID format")
		return
	}

	user, err := h.userService.GetUserByID(userIDStr)
	if err != nil {
		NotFoundError(c, "User not found")
		return
	}

	SuccessResponse(c, "Profile retrieved successfully", user)
}

func (h *UserHandler) SyncUserRoles(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Dapatkan current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		UnauthorizedError(c, "user not found in context")
		return
	}
	currentUserDomain := currentUser.(*domain.User)

	// Dapatkan permissions current user
	currentPermissions, err := h.userService.GetUserPermissions(currentUserDomain.ID)
	if err != nil {
		InternalServerError(c, "failed to get user permissions")
		return
	}

	err = h.userService.SyncUserRoles(userID, req.Roles, currentUserDomain.ID, currentPermissions)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "User roles synced successfully", nil)
}

func (h *UserHandler) SyncUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	// Dapatkan current user dari context
	currentUser, exists := c.Get("user")
	if !exists {
		UnauthorizedError(c, "user not found in context")
		return
	}
	currentUserDomain := currentUser.(*domain.User)

	// Dapatkan permissions current user
	currentPermissions, err := h.userService.GetUserPermissions(currentUserDomain.ID)
	if err != nil {
		InternalServerError(c, "failed to get user permissions")
		return
	}

	err = h.userService.SyncUserPermissions(userID, req.Permissions, currentUserDomain.ID, currentPermissions)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "User permissions synced successfully", nil)
}

func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	userWithPermissions, err := h.userService.GetUserWithRolesAndPermissions(userID)
	if err != nil {
		NotFoundError(c, "User not found")
		return
	}

	SuccessResponse(c, "User permissions retrieved successfully", userWithPermissions)
}
