package payroll

import (
	"net/http"
	"strconv"
	"user_api/internal/common"
	"user_api/internal/dto/pagination"
	"user_api/internal/dto/payroll"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PayrollHandler struct {
	service  *PayrollService
	validate *validator.Validate
}

func NewPayrollHandler(service *PayrollService, validate *validator.Validate) *PayrollHandler {
	return &PayrollHandler{service: service, validate: validate}
}

// GetAllPayrolls godoc
// @Summary List payrolls (paginated)
// @Description Get payrolls with optional pagination query parameters
// @Tags payrolls
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Security BearerAuth
// @Success 200 {object} pagination.PaginationResponse
// @Router /payrolls [get]
func (h *PayrollHandler) GetAllPayrolls(c *gin.Context) {
	req, ok := pagination.ParseFromQuery(c)
	if !ok {
		return
	}

	payrolls, total, err := h.service.GetPayrollsPaginated(req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	c.JSON(http.StatusOK, pagination.BuildResponse(req.Page, req.PageSize, total, payrolls))
}

// CalculatePayroll godoc
// @Summary Calculate payroll
// @Description Calculate net salary for a payroll payload
// @Tags payrolls
// @Accept json
// @Produce json
// @Param payroll body dtos.CreatePayrollDto true "Payroll payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.APIError
// @Security BearerAuth
// @Router /payroll/calculate [post]
func (h *PayrollHandler) CalculatePayroll(c *gin.Context) {
	var dto dtos.CreatePayrollDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	validationErrors, err := common.ValidateStruct(h.validate, dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}
	netSalary := dto.BaseSalary + dto.Bonus - dto.Deductions
	resp := dtos.ResponsePayrollDto{
		WorkerID:   dto.WorkerID,
		Month:      dto.Month,
		BaseSalary: dto.BaseSalary,
		Bonus:      dto.Bonus,
		Deductions: dto.Deductions,
		NetSalary:  netSalary,
	}

	c.JSON(http.StatusOK, resp)
}

// GetPayrollByWorkerId godoc
// @Summary Get payrolls for a worker
// @Description Get payrolls by worker ID
// @Tags payrolls
// @Produce json
// @Param workerId path int true "Worker ID"
// @Security BearerAuth
// @Router /payroll/:workerId [get]
func (h *PayrollHandler) GetPayrollByWorkerId(c *gin.Context) {
	workerId, err := strconv.Atoi(c.Param("workerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid worker ID", Code: "400"})
		return
	}
	payrolls, err := h.service.GetPayrollsByWorkerId(workerId)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Payrolls not found", Code: "404"})
		return
	}
	c.JSON(http.StatusOK, payrolls)
}
