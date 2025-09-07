package handler

import (
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Registration failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	CreatedResponse(c, "User created successfully from public", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid login request", err.Error())
		return
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		UnauthorizedError(c, "Invalid credentials")
		return
	}

	SuccessResponse(c, "Login successful", authResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		UnauthorizedError(c, err.Error())
		return
	}

	//c.JSON(http.StatusOK, authResponse)
	SuccessResponse(c, "Refresh token successful", authResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT setup, logout is handled client-side by removing the token
	// For server-side logout, you might want to implement a token blacklist
	//c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	SuccessResponse(c, "Logged out successfully", nil)
}
