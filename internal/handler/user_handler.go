package handler

import (
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/service"
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

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		NotFoundError(c, err.Error())
		return
	}

	//c.JSON(http.StatusOK, user)
	SuccessResponse(c, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(id)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		InternalServerError(c, err.Error())
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	err := h.userService.ChangePassword(id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		InternalServerError(c, err.Error())
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
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
