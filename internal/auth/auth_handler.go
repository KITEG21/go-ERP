package auth

import (
	"net/http"
	"user_api/internal/common"
	"user_api/internal/dto/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	service  *AuthService
	validate *validator.Validate
	Logger   zerolog.Logger
}

func NewAuthHandler(service *AuthService, validate *validator.Validate, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{service: service, validate: validate, Logger: log}
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
	h.Logger.Info().Msg("Received request to register a new user")
	var dto auth.RegisterDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.Logger.Warn().Err(err).Msg("Failed to bind registration JSON payload")
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)
	if err != nil {
		// This error typically indicates an issue with the validation setup itself, not the input data.
		h.Logger.Error().Err(err).Msg("Error during struct validation setup")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: "Internal validation error", Code: "500"})
		return
	}
	if len(validationErrors) > 0 {
		h.Logger.Warn().Interface("validationErrors", validationErrors).Msg("User registration validation failed")
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	token, err := h.service.Register(dto)
	if err != nil {
		h.Logger.Error().Err(err).Str("email", dto.Email).Msg("Failed to register user in service")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	h.Logger.Info().Str("email", dto.Email).Msg("User registered successfully")
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
	h.Logger.Info().Msg("Received request to login user")
	var dto auth.LoginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.Logger.Warn().Err(err).Msg("Failed to bind login JSON payload")
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Error during struct validation setup for login")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: "Internal validation error", Code: "500"})
		return
	}
	if len(validationErrors) > 0 {
		h.Logger.Warn().Interface("validationErrors", validationErrors).Msg("User login validation failed")
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	token, err := h.service.Login(dto)
	if err != nil {
		h.Logger.Warn().Err(err).Str("email", dto.Email).Msg("Login failed: invalid credentials or service error")
		c.JSON(http.StatusUnauthorized, common.APIError{Message: err.Error(), Code: "401"})
		return
	}

	h.Logger.Info().Str("email", dto.Email).Msg("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}
