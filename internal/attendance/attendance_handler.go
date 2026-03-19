package attendance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"user_api/internal/common"
	"user_api/internal/dto/attendance"
	"user_api/internal/dto/pagination"
)

type AttendanceHandler struct {
	service *AttendanceService
}

func NewAttendanceHandler(service *AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// GetAllAttendances godoc
// @Summary List attendances (paginated)
// @Description Get attendances with optional pagination query parameters
// @Tags attendances
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Security BearerAuth
// @Success 200 {object} pagination.PaginationResponse
// @Router /attendances [get]
func (h *AttendanceHandler) GetAllAttendance(c *gin.Context) {
	req, ok := pagination.ParseFromQuery(c)
	if !ok {
		return
	}

	attendances, total, err := h.service.GetAttendancesPaginated(req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	c.JSON(http.StatusOK, pagination.BuildResponse(req.Page, req.PageSize, total, attendances))
}

// CreateAttendance godoc
// @Summary Create an attendance record
// @Description Create a new attendance record for a worker
// @Tags attendances
// @Accept json
// @Produce json
// @Param attendance body dtos.CreateAttendanceDto true "Attendance payload"
// @Security BearerAuth
// @Router /attendances/checkin [post]
func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	var dto dtos.CreateAttendanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	attendance := Attendance{
		WorkerID: dto.WorkerID,
		CheckIn:  dto.CheckIn,
		CheckOut: dto.CheckOut,
		Date:     dto.Date,
	}

	if err := h.service.CreateAttendance(&attendance); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.JSON(http.StatusCreated, attendance)
}

// UpdateAttendance godoc
// @Summary Update an attendance record
// @Description Update the check-out time of an existing attendance record
// @Tags attendances
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param attendance body dtos.UpdateAttendanceDto true "Updated attendance payload"
// @Security BearerAuth
// @Router /attendances/checkout [put]
func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	var dto dtos.UpdateAttendanceDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}
	attendance, err := h.service.GetAttendanceById(dto.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Attendance not found", Code: "404"})
		return
	}
	checkoutTime, err := time.Parse("15:04:05", *dto.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid time format", Code: "400"})
		return
	}
	attendance.CheckOut = checkoutTime

	if err := h.service.UpdateAttendance(&attendance); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

// GetAttendanceByID godoc
// @Summary Get an attendance record by ID
// @Description Get an attendance record by its ID
// @Tags attendances
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} attendance.Attendance
// @Security BearerAuth
// @Router /attendances/{id} [get]
func (h *AttendanceHandler) GetAttendanceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid ID", Code: "400"})
		return
	}
	attendance, err := h.service.GetAttendanceById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Attendance not found", Code: "404"})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

// GetAttendancesByWorkerID godoc
// @Summary Get attendance records by worker ID
// @Description Get all attendance records for a specific worker
// @Tags attendances
// @Produce json
// @Param worker_id path int true "Worker ID"
// @Success 200 {array} attendance.Attendance
// @Security BearerAuth
// @Router /attendances/worker/{worker_id} [get]
func (h *AttendanceHandler) GetAttendancesByWorkerID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("worker_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid Worker ID", Code: "400"})
		return
	}
	attendances, err := h.service.GetAttendancesByWorkerId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}
	c.JSON(http.StatusOK, attendances)
}
