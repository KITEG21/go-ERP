package payroll

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user_api/internal/common"
	"user_api/internal/dto/payroll"
)

type PayrollHandler struct {
	service *PayrollService
}

func NewPayrollHandler(service *PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// GetAllPayrolls godoc
// @Summary List payrolls
// @Description Get all payrolls
// @Tags payrolls
// @Produce json
// @Router /payrolls [get]
func (h *PayrollHandler) GetAllPayrolls(c *gin.Context) {
	payrolls, _ := h.service.GetAllPayrolls()
	c.JSON(http.StatusOK, payrolls)
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
// @Router /payroll/calculate [post]
func (h *PayrollHandler) CalculatePayroll(c *gin.Context) {
	var dto dtos.CreatePayrollDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
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

	// map[string]interface{}{
	// 	"worker_id":   dto.WorkerID,
	// 	"month":       dto.Month,
	// 	"base_salary": dto.BaseSalary,
	// 	"bonus":       dto.Bonus,
	// 	"deductions":  dto.Deductions,
	// 	"net_salary":  netSalary,
	//}
	c.JSON(http.StatusOK, resp)
}

// GetPayrollByWorkerId godoc
// @Summary Get payrolls for a worker
// @Description Get payrolls by worker ID
// @Tags payrolls
// @Produce json
// @Param workerId path int true "Worker ID"
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
