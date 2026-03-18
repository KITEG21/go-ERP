package auth

import (
	"net/http"
	"user_api/internal/common"
	"user_api/internal/dto/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
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
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	if err := validate.Struct(dto); err != nil {
		var validationErrors []common.ValidationErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, common.ValidationErrorResponse{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Param(),
				Message: common.ValidationErrorResponse{}.CustomErrorMessage(err),
			})
		}
		c.JSON(http.StatusBadRequest, validationErrors)
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
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	err := validate.Struct(dto)
	if err != nil {
		var validationErrors []common.ValidationErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, common.ValidationErrorResponse{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Param(),
				Message: common.ValidationErrorResponse{}.CustomErrorMessage(err),
			})
		}
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	token, err := h.service.Login(dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.APIError{Message: err.Error(), Code: "401"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
