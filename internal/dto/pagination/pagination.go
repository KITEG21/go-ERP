package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"user_api/internal/common"
)

type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// ParseFromQuery extracts and validates pagination parameters from Gin query string.
// Returns default page=1, pageSize=10 if not provided; aborts with 400 if invalid.
func ParseFromQuery(c *gin.Context) (PaginationRequest, bool) {
	req := PaginationRequest{
		Page:     1,
		PageSize: 10,
	}

	if p := c.Query("page"); p != "" {
		parsed, err := strconv.Atoi(p)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid page", Code: "400"})
			return req, false
		}
		req.Page = parsed
	}

	if s := c.Query("page_size"); s != "" {
		parsed, err := strconv.Atoi(s)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, common.APIError{Message: "Invalid page_size", Code: "400"})
			return req, false
		}
		req.PageSize = parsed
	}

	return req, true
}

// BuildResponse constructs a paginated response with calculated total_pages.
func BuildResponse(page int, pageSize int, total int64, data interface{}) PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}
