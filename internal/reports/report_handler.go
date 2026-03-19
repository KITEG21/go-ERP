package reports

import (
	"net/http"
	"strconv"
	"time"
	"user_api/internal/common"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	s *ReportService
}

func NewReportHandler(s *ReportService) *ReportHandler {
	return &ReportHandler{s: s}
}

// GetWorkerAttendanceReport godoc
// @Summary Get worker attendance report
// @Description Get aggregated attendance data for workers within a date range
// @Tags reports
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Param worker_id query int false "Filter by worker ID"
// @Success 200 {array} report.WorkerAttendanceReport
// @Security BearerAuth
// @Router /reports/workers/attendance [get]
func (h *ReportHandler) GetWorkerAttendanceReport(c *gin.Context) {
	// Parse query parameters
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")
	departmentIDStr := c.DefaultQuery("department_id", "0")
	workerIDStr := c.DefaultQuery("worker_id", "0")

	// Validate dates
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, common.APIError{
			Message: "start_date and end_date are required",
			Code:    "400",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{
			Message: "Invalid start_date format (use YYYY-MM-DD)",
			Code:    "400",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{
			Message: "Invalid end_date format (use YYYY-MM-DD)",
			Code:    "400",
		})
		return
	}

	departmentID, _ := strconv.Atoi(departmentIDStr)
	workerID, _ := strconv.Atoi(workerIDStr)

	// Call service
	report, err := h.s.GenerateWorkerAttendanceReport(startDate, endDate, departmentID, workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{
			Message: err.Error(),
			Code:    "500",
		})
		return
	}

	c.JSON(http.StatusOK, report)
}
