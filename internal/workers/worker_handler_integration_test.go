package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"user_api/internal/common"
	"user_api/internal/database"
	"user_api/internal/departments"
)

func setupWorkerRouter(t *testing.T) (*gin.Engine, int) {
	t.Helper()

	_ = godotenv.Load("../../.env.test")

	database.Connect()
	// Optional: migrate only needed tables
	database.DB.AutoMigrate(&Worker{}, &departments.Department{})

	// Clean up any existing records to make the tests deterministic
	database.DB.Exec("DELETE FROM payrolls")
	database.DB.Exec("DELETE FROM attendances")
	database.DB.Exec("DELETE FROM workers")
	database.DB.Exec("DELETE FROM departments")

	// Ensure a department exists for the FK constraint
	testDept := departments.Department{Name: "Test Dept"}
	err := database.DB.Create(&testDept).Error
	require.NoError(t, err)

	repo := &WorkerRepository{}
	svc := NewWorkerService(repo)
	validate := common.NewValidator()
	handler := NewWorkerHandler(svc, validate)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Only worker routes, no auth middleware
	r.POST("/workers", handler.CreateWorker)
	r.GET("/workers/:id", handler.GetWorkerById)
	r.GET("/workers", handler.GetAllWorkers)

	return r, testDept.ID
}

// Happy path: valid worker creation
func TestCreateWorker_Valid(t *testing.T) {
	r, deptID := setupWorkerRouter(t)

	payload := map[string]interface{}{
		"name":          "Test User",
		"email":         "test.user@example.com",
		"department_id": deptID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

// Validation: missing email
func TestCreateWorker_MissingEmail(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	payload := map[string]string{
		"name": "Test User",
		// email missing
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// Validation: invalid email format
func TestCreateWorker_InvalidEmail(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	payload := map[string]string{
		"name":  "Test User",
		"email": "not-an-email", // Invalid email
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// Validation: name too short (min=3)
func TestCreateWorker_NameTooShort(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	payload := map[string]string{
		"name":  "AB", // Only 2 chars, min is 3
		"email": "test@example.com",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// Validation: empty name (notblank fails)
func TestCreateWorker_EmptyName(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	payload := map[string]string{
		"name":  "", // Empty violates notblank
		"email": "test@example.com",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/workers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// Happy path: get existing worker
func TestGetWorkerById_Found(t *testing.T) {
	r, deptID := setupWorkerRouter(t)

	// Insert a worker directly to the DB (setup)
	worker := Worker{Name: "Maria", Email: "maria@example.com", DepartmentId: &deptID}
	err := database.DB.Create(&worker).Error
	require.NoError(t, err)
	defer database.DB.Delete(&worker) // Cleanup

	req := httptest.NewRequest(http.MethodGet, "/workers/"+itoa(worker.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

// Error: get non-existent worker
func TestGetWorkerById_NotFound(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/workers/999999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

// Error: invalid ID format
func TestGetWorkerById_InvalidID(t *testing.T) {
	r, _ := setupWorkerRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/workers/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
func TestGetAllWorkers_Paginated(t *testing.T) {
	r, deptID := setupWorkerRouter(t)

	for i := 1; i <= 15; i++ {
		worker := Worker{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i), DepartmentId: &deptID}
		err := database.DB.Create(&worker).Error
		require.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/workers?page=2&page_size=5", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Page       int      `json:"page"`
		PageSize   int      `json:"page_size"`
		Total      int64    `json:"total"`
		TotalPages int      `json:"total_pages"`
		Data       []Worker `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, 2, resp.Page)
	require.Equal(t, 5, resp.PageSize)
	require.Equal(t, int64(15), resp.Total)
	require.Equal(t, 3, resp.TotalPages)
	require.Len(t, resp.Data, 5)
}

// helper to avoid strconv import in example
func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
