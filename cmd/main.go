// @title User API
// @version 1.0
// @description Simple ERP API
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"os"
	_ "user_api/cmd/docs"
	"user_api/internal/attendance"
	"user_api/internal/auth"
	"user_api/internal/database"
	"user_api/internal/departments"
	"user_api/internal/middleware"
	"user_api/internal/payroll"
	"user_api/internal/reports"
	"user_api/internal/workers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	_ = godotenv.Load("../.env")

	database.Connect()
	if err := database.RunMigrations(); err != nil {
		panic(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret"
	}

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.SetTrustedProxies([]string{"192.168.1.2", "127.168.1.2"})
	jwtService := auth.NewJWTService(jwtSecret)

	workerRepo := &workers.WorkerRepository{}
	workerService := workers.NewWorkerService(workerRepo)
	workerHandler := workers.NewWorkerHandler(workerService)

	departmentRepo := departments.DepartmentRepository{}
	departmentService := departments.NewDepartmentService(departmentRepo)
	departmentHandler := departments.NewDepartmentHandler(departmentService)

	attendanceRepo := attendance.AttendanceRepository{}
	attendanceService := attendance.NewAttendanceService(&attendanceRepo)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceService)

	payrollRepo := payroll.PayrollRepository{}
	payrollService := payroll.NewPayrollService(payrollRepo)
	payrollHandler := payroll.NewPayrollHandler(payrollService)

	reportService := reports.ReportService{}
	reportHandler := reports.NewReportHandler(&reportService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRoutes := r.Group("/api/v1/auth")
	{
		userRepo := auth.NewUserRepository()
		authService := auth.NewAuthService(userRepo, jwtService)
		authHandler := auth.NewAuthHandler(authService)

		publicRoutes.POST("/register", authHandler.Register)
		publicRoutes.POST("/login", authHandler.Login)
	}

	r.Use(middleware.AuthMiddleware(jwtService))
	api := r.Group("/api/v1")

	//Workers
	api.GET("/", workerHandler.TestHandler)
	api.POST("/workers", workerHandler.CreateWorker)
	api.GET("/workers", workerHandler.GetAllWorkers)
	api.GET("/workers/:id", workerHandler.GetWorkerById)
	api.PUT("/workers", workerHandler.UpdateWorker)
	api.DELETE("/workers/:id", workerHandler.DeleteWorker)

	//Dptos
	api.GET("/departments", departmentHandler.GetAllDepartments)
	api.POST("/departments", departmentHandler.CreateDepartment)
	api.GET("/departments/:id", departmentHandler.GetDepartmentByID)
	api.PUT("/departments", departmentHandler.UpdateDepartment)
	api.DELETE("/departments/:id", departmentHandler.DeleteDepartment)

	//Attendance
	api.GET("/attendances", attendanceHandler.GetAllAttendance)
	api.POST("/attendances", attendanceHandler.CreateAttendance)
	api.PUT("/attendances/:id", attendanceHandler.UpdateAttendance)
	api.GET("/attendances/:id", attendanceHandler.GetAttendanceByID)
	api.GET("/attendances/worker/:worker_id", attendanceHandler.GetAttendancesByWorkerID)

	api.GET("/payrolls", payrollHandler.GetAllPayrolls)
	api.POST("/payroll/calculate", payrollHandler.CalculatePayroll)
	api.GET("/payrolls/:workerId", payrollHandler.GetPayrollByWorkerId)

	api.GET("/reports/workers/attendance", reportHandler.GetWorkerAttendanceReport)

	r.Run(":8080")

}
