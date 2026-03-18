package payroll

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

	"user_api/internal/database"
	"user_api/internal/departments"
	"user_api/internal/workers"
)

func setupPayrollRouter(t *testing.T) (*gin.Engine, int) {
	t.Helper()

	_ = godotenv.Load("../../.env")

	database.Connect()
	database.DB.AutoMigrate(
		&departments.Department{},
		&workers.Worker{},
		&Payroll{},
	)

	dept := departments.Department{Name: "Test Dept"}
	require.NoError(t, database.DB.Create(&dept).Error)
	t.Cleanup(func() { database.DB.Delete(&dept) })

	worker := workers.Worker{Name: "Test Worker", Email: "test@example.com", DepartmentId: &dept.ID}
	require.NoError(t, database.DB.Create(&worker).Error)
	t.Cleanup(func() { database.DB.Delete(&worker) })

	repo := &PayrollRepository{}
	svc := NewPayrollService(*repo)
	handler := NewPayrollHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/payrolls", handler.GetAllPayrolls)
	r.POST("/payroll/calculate", handler.CalculatePayroll)
	r.GET("/payrolls/:workerId", handler.GetPayrollByWorkerId)

	return r, worker.ID
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func TestCalculatePayroll_Valid(t *testing.T) {
	r, workerID := setupPayrollRouter(t)

	payload := map[string]interface{}{
		"worker_id":   workerID,
		"month":       "2026-03",
		"base_salary": 1000.0,
		"bonus":       100.0,
		"deductions":  50.0,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/payroll/calculate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, float64(1050), resp["net_salary"]) // JSON numbers decode to float64
}

func TestCalculatePayroll_MissingMonth(t *testing.T) {
	r, workerID := setupPayrollRouter(t)

	payload := map[string]interface{}{
		"worker_id":   workerID,
		"base_salary": 1000.0,
		"bonus":       100.0,
		"deductions":  50.0,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/payroll/calculate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPayrollsByWorkerId_NoRecords(t *testing.T) {
	r, workerID := setupPayrollRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/payrolls/"+itoa(workerID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var payrolls []interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payrolls))
	require.Len(t, payrolls, 0)
}

func TestGetPayrollsByWorkerId_Found(t *testing.T) {
	r, workerID := setupPayrollRouter(t)

	record := Payroll{
		WorkerId:   workerID,
		Month:      "2026-03",
		BaseSalary: 1000,
		Bonus:      100,
		Deductions: 50,
		NetSalary:  1050,
		Status:     Processed,
	}
	require.NoError(t, database.DB.Create(&record).Error)
	t.Cleanup(func() { database.DB.Delete(&record) })

	req := httptest.NewRequest(http.MethodGet, "/payrolls/"+itoa(workerID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var payrolls []Payroll
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payrolls))
	require.Len(t, payrolls, 1)
	require.Equal(t, float32(1050), payrolls[0].NetSalary)
}
