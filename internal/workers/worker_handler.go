package workers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"user_api/internal/common"
	"user_api/internal/dto/worker"
)

type WorkerHandler struct {
	service *WorkerService
}

func NewWorkerHandler(service *WorkerService) *WorkerHandler {
	return &WorkerHandler{service: service}
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
// @Summary List workers
// @Description Get all workers
// @Tags workers
// @Produce json
// @Success 200 {array} workers.Worker
// @Router /workers [get]
func (h *WorkerHandler) GetAllWorkers(c *gin.Context) {
	workers, _ := h.service.GetAllWorkers()
	c.JSON(http.StatusOK, workers)
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
// @Router /workers [post]
func (h *WorkerHandler) CreateWorker(c *gin.Context) {
	var dto dtos.CreateWorkerDto
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, common.APIError{Message: err.Error(), Code: "401"})
		return
	}
	err := validate.Struct(dto)

	if err != nil {
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
	var worker = Worker{
		Name:  dto.Name,
		Email: dto.Email,
	}
	h.service.CreateWorker(&worker)
	c.JSON(http.StatusCreated, dto)
}

// GetWorkerById godoc
// @Summary Get a worker
// @Description Get a worker by ID
// @Tags workers
// @Produce json
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
// @Router /workers [put]
func (h *WorkerHandler) UpdateWorker(c *gin.Context) {
	var dto dtos.UpdateWorkerDto
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

	worker, err := h.service.GetWorkerById(dto.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIError{Message: "Worker not found", Code: "404"})
		return
	}
	worker.Name = dto.Name
	worker.Email = dto.Email
	h.service.UpdateWorker(&worker)
	c.JSON(http.StatusOK, worker)
}

// DeleteWorker godoc
// @Summary Delete a worker
// @Description Delete a worker by ID
// @Router /workers/{id} [delete]
// @Tags workers
func (h *WorkerHandler) DeleteWorker(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.service.DeleteWorker(id)
	c.Status(http.StatusNoContent)
}
