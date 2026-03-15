package report

import "time"

type ReportFilters struct {
	StartDate    time.Time `json:"startDate" example:"2024-01-01"`
	EndDate      time.Time `json:"endDate" example:"2024-01-31"`
	DepartmentId int       `json:"departmentId" example:"1"`
	WorkerId     int       `json:"workerId" example:"1"`
}
