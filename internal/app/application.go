package app

import (
	"github.com/gin-gonic/gin"
	docs "user_api/cmd/docs"
	"user_api/internal/attendance"
	"user_api/internal/auth"
	"user_api/internal/departments"
	"user_api/internal/middleware"
	"user_api/internal/payroll"
	"user_api/internal/reports"
	"user_api/internal/workers"
)

type Application struct {
	Engine *gin.Engine
}

func NewApp(jwtSecret string) (*Application, error) {

	jwtService := auth.NewJWTService(jwtSecret)

	workerRepo := &workers.WorkerRepository{}
	workerService := workers.NewWorkerService(workerRepo)
	workerHandler := workers.NewWorkerHandler(workerService)

	deptRepo := departments.DepartmentRepository{}
	deptSvc := departments.NewDepartmentService(deptRepo)
	deptHandler := departments.NewDepartmentHandler(deptSvc)

	attendanceRepo := attendance.AttendanceRepository{}
	attendanceService := attendance.NewAttendanceService(&attendanceRepo)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceService)

	payrollRepo := payroll.PayrollRepository{}
	payrollService := payroll.NewPayrollService(payrollRepo)
	payrollHandler := payroll.NewPayrollHandler(payrollService)

	reportService := reports.ReportService{}
	reportHandler := reports.NewReportHandler(&reportService)

	userRepo := auth.NewUserRepository()
	authService := auth.NewAuthService(userRepo, jwtService)
	authHandler := auth.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Ensure documentation paths use API version prefix from routing
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/doc.json", func(c *gin.Context) {
		c.Data(200, "application/json; charset=utf-8", []byte(docs.SwaggerInfo.ReadDoc()))
	})
	r.GET("/scalar", func(c *gin.Context) {
		html := `
		<!doctype html>
		<html>
		  <head>
			<title>API Documentation</title>
			<meta charset="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
		  </head>
		  <body>
			<script id="api-reference" data-url="/swagger/doc.json"></script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		  </body>
		</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	registerAuthRoutes(r, authHandler)
	registerAPIRoutes(r, workerHandler, deptHandler, attendanceHandler, payrollHandler, reportHandler, jwtService)

	return &Application{Engine: r}, nil
}
