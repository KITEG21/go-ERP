package auth

import (
	"net/http"
	"user_api/internal/common"
	"user_api/internal/dto/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterDto true "User registration data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.APIError
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto auth.RegisterDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	token, err := h.service.Register(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginDto true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.APIError
// @Failure 401 {object} common.APIError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto auth.LoginDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	token, err := h.service.Login(dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.APIError{Message: err.Error(), Code: "401"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
