package handler

import (
	"fmt"
	"os"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

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

	// Pastikan tidak ada role IDs yang dikirim melalui register publik
	req.RoleIDs = nil

	user, err := h.authService.Register(req)
	if err != nil {
		HandleError(c, err)
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
		HandleError(c, err)
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
		HandleError(c, err)
		return
	}

	//c.JSON(http.StatusOK, authResponse)
	SuccessResponse(c, "Refresh token successful", authResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Ambil userID dari context (diset oleh AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		UnauthorizedError(c, "User ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		InternalServerError(c, "Invalid user ID format")
		return
	}

	// Panggil service untuk logout
	err := h.authService.Logout(userIDStr)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Logged out successfully", nil)
}

func (h *AuthHandler) ServeFile(c *gin.Context) {
	folder := c.Param("folder")     // e.g., "students"
	filename := c.Param("filename") // e.g., "akta_xyz.pdf"

	// Validasi folder agar user tidak bisa akses folder sistem (Path Traversal Attack)
	if folder != "students" && folder != "employees" {
		c.JSON(403, gin.H{"error": "Forbidden access"})
		return
	}

	targetPath := fmt.Sprintf("./storage/uploads/%s/%s", folder, filename)

	// Cek file ada atau tidak
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	// Sajikan file
	c.File(targetPath)
}
