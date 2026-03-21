package departments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"user_api/internal/common"
	"user_api/internal/database"
)

func setupDepartmentRouter(t *testing.T) *gin.Engine {
	t.Helper()

	_ = godotenv.Load("../../.env.test")

	database.Connect()
	database.DB.AutoMigrate(&Department{})

	log := zerolog.New(zerolog.NewTestWriter(t))
	repo := &DepartmentRepository{}
	svc := NewDepartmentService(repo, log)
	validate := common.NewValidator()
	handler := NewDepartmentHandler(svc, validate, log)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/departments", handler.CreateDepartment)
	r.GET("/departments/:id", handler.GetDepartmentByID)
	r.GET("/departments", handler.GetAllDepartments)
	r.PUT("/departments", handler.UpdateDepartment)
	r.DELETE("/departments/:id", handler.DeleteDepartment)

	return r
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func TestCreateDepartment_Valid(t *testing.T) {
	r := setupDepartmentRouter(t)

	payload := map[string]string{
		"name":        "Test Dept",
		"description": "A description",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateDepartment_MissingName(t *testing.T) {
	r := setupDepartmentRouter(t)

	payload := map[string]string{
		"description": "No name",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDepartmentById_Found(t *testing.T) {
	r := setupDepartmentRouter(t)

	dept := Department{Name: "Found Dept", Description: "x"}
	err := database.DB.Create(&dept).Error
	require.NoError(t, err)
	defer database.DB.Delete(&dept)

	req := httptest.NewRequest(http.MethodGet, "/departments/"+itoa(dept.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetDepartmentById_NotFound(t *testing.T) {
	r := setupDepartmentRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/departments/999999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetDepartmentById_InvalidID(t *testing.T) {
	r := setupDepartmentRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/departments/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDepartment_Valid(t *testing.T) {
	r := setupDepartmentRouter(t)

	dept := Department{Name: "Original", Description: "old"}
	err := database.DB.Create(&dept).Error
	require.NoError(t, err)
	defer database.DB.Delete(&dept)

	payload := map[string]interface{}{
		"id":          dept.ID,
		"name":        "Updated Name",
		"description": "new",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteDepartment_Valid(t *testing.T) {
	r := setupDepartmentRouter(t)

	dept := Department{Name: "To Delete", Description: "x"}
	err := database.DB.Create(&dept).Error
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/departments/"+itoa(dept.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)

	// Verify it's gone
	var found Department
	err = database.DB.First(&found, dept.ID).Error
	require.Error(t, err)
}
