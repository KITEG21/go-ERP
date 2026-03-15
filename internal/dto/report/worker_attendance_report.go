package report

type WorkerAttendanceReport struct {
	WorkerId       int     `json:"workerId"`
	WorkerName     string  `json:"workerName"`
	Department     string  `json:"department"`
	DaysPresent    int     `json:"daysPresent"`
	DaysAbsent     int     `json:"daysAbsent"`
	HoursWorked    float64 `json:"hoursWorked"`
	AttendanceRate float64 `json:"attendanceRate"`
}
