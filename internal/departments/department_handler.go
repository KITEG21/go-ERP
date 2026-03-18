package departments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"user_api/internal/common"
	"user_api/internal/dto/department"
)

type DepartmentHandler struct {
	service *DepartmentService
}

func NewDepartmentHandler(service *DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// GetAllDepartments godoc
// @Summary List departments
// @Description Get all departments
// @Tags departments
// @Produce json
// @Success 200 {array} departments.Department
// @Router /departments [get]
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, _ := h.service.GetAllDepartments()
	c.JSON(http.StatusOK, departments)
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
// @Router /departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var dto dtos.CreateDepartmentDto
	validate := validator.New()

	validate.RegisterValidation("notblank", validators.NotBlank)
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	if err := validate.Struct(dto); err != nil {
		var validationErrors []common.ValidationErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			var error common.ValidationErrorResponse
			error.Field = err.Field()
			error.Tag = err.Tag()
			error.Value = err.Param()
			error.Message = error.CustomErrorMessage(err)
			validationErrors = append(validationErrors, error)
		}
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	department := &Department{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := h.service.CreateDepartment(department); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	c.JSON(http.StatusCreated, department)
}

// GetDepartmentById godoc
// @Summary Get a department
// @Description Get a department by ID
// @Tags departments
// @Produce json
// @Success 200 {object} departments.Department
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")
	departmentId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid ID", Code: "400"})
		return
	}
	department, err := h.service.GetDepartmentById(departmentId)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Department not found", Code: "404"})
		return
	}
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
// @Router /departments [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	var dto dtos.UpdateDepartmentDto
	validate := validator.New()

	validate.RegisterValidation("notblank", validators.NotBlank)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	if err := validate.Struct(dto); err != nil {
		var validationErrors []common.ValidationErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			var error common.ValidationErrorResponse
			error.Field = err.Field()
			error.Tag = err.Tag()
			error.Value = err.Param()
			error.Message = error.CustomErrorMessage(err)
			validationErrors = append(validationErrors, error)
		}
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	department, err := h.service.GetDepartmentById(dto.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Department not found", Code: "404"})
		return
	}
	department.Name = dto.Name
	department.Description = dto.Description
	if err := h.service.UpdateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.JSON(http.StatusOK, department)
}

// DeleteDepartment godoc
// @Summary Delete a department
// @Description Delete a department by ID
// @Router /departments/{id} [delete]
// @Tags departments
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	departmentId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid ID", Code: "400"})
		return
	}
	if err := h.service.DeleteDepartment(departmentId); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.Status(http.StatusNoContent)
}
