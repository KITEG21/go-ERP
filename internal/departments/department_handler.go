package departments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"user_api/internal/common"
	dtos "user_api/internal/dto/department"
	"user_api/internal/dto/pagination"
)

type DepartmentHandler struct {
	service  *DepartmentService
	validate *validator.Validate
	Logger   zerolog.Logger
}

func NewDepartmentHandler(service *DepartmentService, validate *validator.Validate, log zerolog.Logger) *DepartmentHandler {
	return &DepartmentHandler{service: service, validate: validate, Logger: log}
}

// GetAllDepartments godoc
// @Summary List departments (paginated)
// @Description Get departments with optional pagination query parameters
// @Tags departments
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Security BearerAuth
// @Success 200 {object} pagination.PaginationResponse
// @Router /departments [get]
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	h.Logger.Info().Msg("Received request to list all departments")
	req, ok := pagination.ParseFromQuery(c)
	if !ok {
		h.Logger.Warn().Msg("Failed to parse pagination query parameters for departments")
		return // Error already handled by ParseFromQuery returning false
	}

	departments, total, err := h.service.GetDepartmentsPaginated(req.Page, req.PageSize)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to retrieve paginated departments from service")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	h.Logger.Info().Int("count", len(departments)).Msg("Successfully retrieved departments")
	c.JSON(http.StatusOK, pagination.BuildResponse(req.Page, req.PageSize, total, departments))
}

// CreateDepartment godoc
// @Summary Create a department
// @Description Create a new department
// @Tags departments
// @Accept json
// @Produce json
// @Param department body dtos.CreateDepartmentDto true "Department payload"
// @Success 201 {object} dtos.CreateDepartmentDto
// @Failure 400 {object} common.APIError
// @Security BearerAuth
// @Router /departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	h.Logger.Info().Msg("Received request to create a new department")
	var dto dtos.CreateDepartmentDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.Logger.Warn().Err(err).Msg("Failed to bind create department JSON payload")
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Error during department struct validation setup")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: "Internal validation error", Code: "500"})
		return
	}
	if len(validationErrors) > 0 {
		h.Logger.Warn().Interface("validationErrors", validationErrors).Msg("Department creation validation failed")
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	department := &Department{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := h.service.CreateDepartment(department); err != nil {
		h.Logger.Error().Err(err).Str("departmentName", dto.Name).Msg("Failed to create department in service")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	h.Logger.Info().Str("departmentName", dto.Name).Msg("Department created successfully")
	c.JSON(http.StatusCreated, department)
}

// GetDepartmentById godoc
// @Summary Get a department
// @Description Get a department by ID
// @Tags departments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} departments.Department
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")
	h.Logger.Info().Str("idParam", id).Msg("Received request to get department by ID")

	departmentId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Warn().Str("idParam", id).Msg("Invalid ID format for get department request")
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid ID", Code: "400"})
		return
	}

	department, err := h.service.GetDepartmentById(departmentId)
	if err != nil {
		h.Logger.Warn().Int("departmentID", departmentId).Msg("Department not found")
		c.JSON(http.StatusNotFound, common.APIError{Message: "Department not found", Code: "404"})
		return
	}

	h.Logger.Info().Int("departmentID", departmentId).Msg("Successfully retrieved department")
	c.JSON(http.StatusOK, department)
}

// UpdateDepartment godoc
// @Summary Update a department
// @Description Update an existing department
// @Tags departments
// @Accept json
// @Produce json
// @Param department body dtos.UpdateDepartmentDto true "Department payload"
// @Success 201 {object} dtos.UpdateDepartmentDto
// @Failure 400 {object} common.APIError
// @Security BearerAuth
// @Router /departments [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	h.Logger.Info().Msg("Received request to update a department")
	var dto dtos.UpdateDepartmentDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		h.Logger.Warn().Err(err).Msg("Failed to bind update department JSON payload")
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Error during department update struct validation setup")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: "Internal validation error", Code: "500"})
		return
	}
	if len(validationErrors) > 0 {
		h.Logger.Warn().Interface("validationErrors", validationErrors).Msg("Department update validation failed")
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	department, err := h.service.GetDepartmentById(dto.Id)
	if err != nil {
		h.Logger.Warn().Int("departmentID", dto.Id).Msg("Department not found for update")
		c.JSON(http.StatusNotFound, common.APIError{Message: "Department not found", Code: "404"})
		return
	}
	department.Name = dto.Name
	department.Description = dto.Description

	if err := h.service.UpdateDepartment(&department); err != nil {
		h.Logger.Error().Err(err).Int("departmentID", dto.Id).Msg("Failed to update department in service")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	h.Logger.Info().Int("departmentID", dto.Id).Msg("Department updated successfully")
	c.JSON(http.StatusOK, department)
}

// DeleteDepartment godoc
// @Summary Delete a department
// @Description Delete a department by ID
// @Router /departments/{id} [delete]
// @Tags departments
// @Security BearerAuth
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	h.Logger.Info().Str("idParam", id).Msg("Received request to delete department by ID")

	departmentId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Warn().Str("idParam", id).Msg("Invalid ID format for delete department request")
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid ID", Code: "400"})
		return
	}

	if err := h.service.DeleteDepartment(departmentId); err != nil {
		h.Logger.Error().Err(err).Int("departmentID", departmentId).Msg("Failed to delete department in service")
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	h.Logger.Info().Int("departmentID", departmentId).Msg("Department deleted successfully")
	c.Status(http.StatusNoContent)
}
