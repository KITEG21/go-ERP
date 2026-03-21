package attendance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"user_api/internal/database"
	"user_api/internal/departments"
	"user_api/internal/workers"
)

func setupAttendanceRouter(t *testing.T) (*gin.Engine, int) {
	t.Helper()

	_ = godotenv.Load("../../.env.test")

	database.Connect()
	database.DB.AutoMigrate(
		&departments.Department{},
		&workers.Worker{},
		&Attendance{},
	)

	// Create a department + worker for FK constraints
	dept := departments.Department{Name: "Test Dept"}
	require.NoError(t, database.DB.Create(&dept).Error)

	worker := workers.Worker{Name: "Test Worker", Email: "test@example.com", DepartmentId: &dept.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	log := zerolog.New(zerolog.NewTestWriter(t))
	repo := &AttendanceRepository{}
	svc := NewAttendanceService(repo, log)
	handler := NewAttendanceHandler(svc, log)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/attendances", handler.CreateAttendance)
	r.PUT("/attendances/:id", handler.UpdateAttendance)
	r.GET("/attendances/:id", handler.GetAttendanceByID)
	r.GET("/attendances/worker/:worker_id", handler.GetAttendancesByWorkerID)
	r.GET("/attendances", handler.GetAllAttendance)

	return r, worker.ID
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func TestCreateAttendance_Valid(t *testing.T) {
	r, workerID := setupAttendanceRouter(t)

	payload := map[string]interface{}{
		"worker_id": workerID,
		"check_in":  time.Now().UTC(),
		"check_out": time.Now().Add(8 * time.Hour).UTC(),
		"date":      time.Now().UTC(),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/attendances", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateAttendance_MissingWorkerID(t *testing.T) {
	r, _ := setupAttendanceRouter(t)

	payload := map[string]interface{}{
		"check_in":  time.Now().UTC(),
		"check_out": time.Now().Add(8 * time.Hour).UTC(),
		"date":      time.Now().UTC(),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/attendances", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAttendanceById_Found(t *testing.T) {
	r, workerID := setupAttendanceRouter(t)

	record := Attendance{
		WorkerID: workerID,
		CheckIn:  time.Now().UTC(),
		CheckOut: time.Now().Add(8 * time.Hour).UTC(),
		Date:     time.Now().UTC(),
	}
	require.NoError(t, database.DB.Create(&record).Error)
	defer database.DB.Delete(&record)

	req := httptest.NewRequest(http.MethodGet, "/attendances/"+itoa(record.ID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAttendanceById_NotFound(t *testing.T) {
	r, _ := setupAttendanceRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/attendances/999999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAttendanceById_InvalidID(t *testing.T) {
	r, _ := setupAttendanceRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/attendances/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAttendance_Valid(t *testing.T) {
	r, workerID := setupAttendanceRouter(t)

	record := Attendance{
		WorkerID: workerID,
		CheckIn:  time.Now().UTC(),
		CheckOut: time.Now().Add(8 * time.Hour).UTC(),
		Date:     time.Now().UTC(),
	}
	require.NoError(t, database.DB.Create(&record).Error)
	defer database.DB.Delete(&record)

	payload := map[string]interface{}{
		"id":        record.ID,
		"check_out": "23:59:59",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/attendances/"+itoa(record.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAttendance_InvalidTimeFormat(t *testing.T) {
	r, workerID := setupAttendanceRouter(t)

	record := Attendance{
		WorkerID: workerID,
		CheckIn:  time.Now().UTC(),
		CheckOut: time.Now().Add(8 * time.Hour).UTC(),
		Date:     time.Now().UTC(),
	}
	require.NoError(t, database.DB.Create(&record).Error)
	defer database.DB.Delete(&record)

	payload := map[string]interface{}{
		"id":        record.ID,
		"check_out": "not-a-time",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/attendances/"+itoa(record.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAttendancesByWorkerID_Found(t *testing.T) {
	r, workerID := setupAttendanceRouter(t)

	record := Attendance{
		WorkerID: workerID,
		CheckIn:  time.Now().UTC(),
		CheckOut: time.Now().Add(8 * time.Hour).UTC(),
		Date:     time.Now().UTC(),
	}
	require.NoError(t, database.DB.Create(&record).Error)
	defer database.DB.Delete(&record)

	req := httptest.NewRequest(http.MethodGet, "/attendances/worker/"+itoa(workerID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAttendancesByWorkerID_NotFound(t *testing.T) {
	r, _ := setupAttendanceRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/attendances/worker/999999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code) // returns [] on no matches
}
