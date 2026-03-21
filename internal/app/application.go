package app

import (
	"time"
	docs "user_api/cmd/docs"
	"user_api/internal/attendance"
	"user_api/internal/auth"
	"user_api/internal/common"
	"user_api/internal/departments"
	"user_api/internal/middleware"
	"user_api/internal/payroll"
	"user_api/internal/reports"
	"user_api/internal/workers"

	tollbooth "github.com/didip/tollbooth/v7"
	tollbooth_gin "github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Application struct {
	Engine *gin.Engine
	Logger zerolog.Logger
}

func NewApp(jwtSecret string, log zerolog.Logger) (*Application, error) {

	jwtService := auth.NewJWTService(jwtSecret)
	validate := common.NewValidator()

	workerRepo := &workers.WorkerRepository{}
	workerService := workers.NewWorkerService(workerRepo, log)
	workerHandler := workers.NewWorkerHandler(workerService, validate, log)

	deptRepo := departments.DepartmentRepository{}
	deptSvc := departments.NewDepartmentService(&deptRepo, log)
	deptHandler := departments.NewDepartmentHandler(deptSvc, validate, log)

	attendanceRepo := attendance.AttendanceRepository{}
	attendanceService := attendance.NewAttendanceService(&attendanceRepo, log)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceService, log)

	payrollRepo := payroll.PayrollRepository{}
	payrollService := payroll.NewPayrollService(&payrollRepo, log)
	payrollHandler := payroll.NewPayrollHandler(payrollService, validate, log)

	reportService := reports.ReportService{}
	reportHandler := reports.NewReportHandler(&reportService, log)

	userRepo := auth.NewUserRepository()
	authService := auth.NewAuthService(userRepo, jwtService, log)
	authHandler := auth.NewAuthHandler(authService, validate, log)

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware(log))
	r.Use(middleware.MetricsMiddleware())
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	limiter := tollbooth.NewLimiter(5, nil) // 5 req/sec default
	limiter.SetBurst(10)
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

	r.Use(tollbooth_gin.LimitHandler(limiter))

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

	return &Application{Engine: r, Logger: log}, nil
}
