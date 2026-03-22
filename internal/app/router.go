package app

import (
	"user_api/internal/attendance"
	"user_api/internal/auth"
	"user_api/internal/departments"
	"user_api/internal/health"
	"user_api/internal/middleware"
	"user_api/internal/payroll"
	"user_api/internal/reports"
	"user_api/internal/workers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerAuthRoutes(r *gin.Engine, h *auth.AuthHandler) {
	public := r.Group("/api/v1/auth")
	public.POST("/register", h.Register)
	public.POST("/login", h.Login)
}

func registerPublicRoutes(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	health.RegisterRoutes(r)
}

func registerAPIRoutes(r *gin.Engine, wh *workers.WorkerHandler, dh *departments.DepartmentHandler, ah *attendance.AttendanceHandler,
	ph *payroll.PayrollHandler, rh *reports.ReportHandler, jwtSvc *auth.JWTService) {
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSvc))

	//Workers
	api.GET("/", wh.TestHandler)
	api.POST("/workers", wh.CreateWorker)
	api.GET("/workers", wh.GetAllWorkers)
	api.GET("/workers/:id", wh.GetWorkerById)
	api.PUT("/workers", wh.UpdateWorker)
	api.DELETE("/workers/:id", wh.DeleteWorker)

	//Dptos
	api.GET("/departments", dh.GetAllDepartments)
	api.POST("/departments", dh.CreateDepartment)
	api.GET("/departments/:id", dh.GetDepartmentByID)
	api.PUT("/departments", dh.UpdateDepartment)
	api.DELETE("/departments/:id", dh.DeleteDepartment)

	//Attendance
	api.GET("/attendances", ah.GetAllAttendance)
	api.POST("/attendances", ah.CreateAttendance)
	api.PUT("/attendances/:id", ah.UpdateAttendance)
	api.GET("/attendances/:id", ah.GetAttendanceByID)
	api.GET("/attendances/worker/:worker_id", ah.GetAttendancesByWorkerID)

	api.GET("/payrolls", ph.GetAllPayrolls)
	api.POST("/payroll/calculate", ph.CalculatePayroll)
	api.GET("/payrolls/:workerId", ph.GetPayrollByWorkerId)

	api.GET("/reports/workers/attendance", rh.GetWorkerAttendanceReport)

}
