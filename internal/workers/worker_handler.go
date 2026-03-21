package workers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"user_api/internal/common"
	"user_api/internal/dto/pagination"
	dtos "user_api/internal/dto/worker"
)

type WorkerHandler struct {
	service  *WorkerService
	validate *validator.Validate
	Logger   zerolog.Logger
}

func NewWorkerHandler(service *WorkerService, validate *validator.Validate, log zerolog.Logger) *WorkerHandler {
	return &WorkerHandler{service: service, validate: validate, Logger: log}
}

// TestHandler godoc
// @Summary Health check
// @Description Returns service status
// @Tags health
// @Produce json
// @Success 200 {string} string "We are up!"
// @Router / [get]
func (h *WorkerHandler) TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "We are up!")
}

// GetAllWorkers godoc
// @Summary List workers (paginated)
// @Description Get workers with optional pagination query parameters
// @Tags workers
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Security BearerAuth
// @Success 200 {object} pagination.PaginationResponse
// @Router /workers [get]
func (h *WorkerHandler) GetAllWorkers(c *gin.Context) {
	req, ok := pagination.ParseFromQuery(c)
	if !ok {
		return
	}

	workers, total, err := h.service.GetWorkersPaginated(req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	c.JSON(http.StatusOK, pagination.BuildResponse(req.Page, req.PageSize, total, workers))
}

// CreateWorker godoc
// @Summary Create a worker
// @Description Create a new worker
// @Tags workers
// @Accept json
// @Produce json
// @Param worker body dtos.CreateWorkerDto true "Worker payload"
// @Success 201 {object} dtos.CreateWorkerDto
// @Failure 400 {object} common.APIError
// @Security BearerAuth
// @Router /workers [post]
func (h *WorkerHandler) CreateWorker(c *gin.Context) {
	var dto dtos.CreateWorkerDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}
	var hireDate string
	if dto.HireDate != nil && *dto.HireDate != "" {
		_, err := time.Parse("2006-01-02", *dto.HireDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid hire_date format, use YYYY-MM-DD", Code: "400"})
			return
		}
		hireDate = *dto.HireDate
	} else {
		hireDate = time.Now().Format("2006-01-02")
	}

	worker := Worker{
		Name:     dto.Name,
		Email:    dto.Email,
		HireDate: hireDate,
	}
	if dto.DepartmentId != nil {
		worker.DepartmentId = dto.DepartmentId
	}

	err = h.service.CreateWorker(&worker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIError{Message: err.Error(), Code: "500"})
		return
	}

	c.JSON(http.StatusCreated, worker)
}

// GetWorkerById godoc
// @Summary Get a worker
// @Description Get a worker by ID
// @Tags workers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} workers.Worker
// @Router /workers/{id} [get]
func (h *WorkerHandler) GetWorkerById(c *gin.Context) {
	id := c.Param("id")
	workerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid worker ID", Code: "400"})
		return
	}
	worker, err := h.service.GetWorkerById(workerId)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Worker not found", Code: "404"})
		return
	}
	c.JSON(http.StatusOK, worker)
}

// UpdateWorker godoc
// @Summary Update a worker
// @Description Update an existing worker
// @Tags workers
// @Accept json
// @Produce json
// @Param worker body dtos.UpdateWorkerDto true "Worker payload"
// @Success 201 {object} dtos.UpdateWorkerDto
// @Failure 400 {object} common.APIError
// @Security BearerAuth
// @Router /workers [put]
func (h *WorkerHandler) UpdateWorker(c *gin.Context) {
	var dto dtos.UpdateWorkerDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "400"})
		return
	}

	validationErrors, err := common.ValidateStruct(h.validate, dto)

	if err != nil {
		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	worker, err := h.service.GetWorkerById(dto.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Worker not found", Code: "404"})
		return
	}
	worker.Name = dto.Name
	worker.Email = dto.Email
	if dto.DepartmentId != nil {
		worker.DepartmentId = dto.DepartmentId
	}
	h.service.UpdateWorker(&worker)
	c.JSON(http.StatusOK, worker)
}

// DeleteWorker godoc
// @Summary Delete a worker
// @Description Delete a worker by ID
// @Router /workers/{id} [delete]
// @Tags workers
// @Security BearerAuth
func (h *WorkerHandler) DeleteWorker(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.service.DeleteWorker(id)
	c.Status(http.StatusNoContent)
}
