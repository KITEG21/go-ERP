package reports

import (
	"fmt"
	"time"
	"user_api/internal/attendance"
	"user_api/internal/database"
	"user_api/internal/dto/report"
)

type ReportService struct{}

func (s *ReportService) GenerateWorkerAttendanceReport(
	startDate, endDate time.Time, deparmentId, workerId int,
) ([]report.WorkerAttendanceReport, error) {

	var atendances []attendance.Attendance
	query := database.DB.
		Joins("JOIN workers ON attendances.worker_id = workers.id").
		Joins("JOIN departments ON workers.department_id = departments.id").
		Preload("Worker.Department")

	if deparmentId > 0 {
		query = query.Where("departments.id = ?", deparmentId)
	}
	if workerId > 0 {
		query = query.Where("workers.id = ?", workerId)
	}

	query = query.
		Where("DATE(attendances.date) BETWEEN ? AND ?", startDate, endDate).
		Order("attendances.worker_id, attendances.date")

	if err := query.Find(&atendances).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch attendance records: %w", err)
	}

	reportMap := make(map[int]*report.WorkerAttendanceReport)

	for _, a := range atendances {
		if _, exists := reportMap[a.WorkerID]; !exists {
			reportMap[a.WorkerID] = &report.WorkerAttendanceReport{
				WorkerId:    a.WorkerID,
				WorkerName:  a.Worker.Name,
				Department:  a.Worker.Department.Name,
				DaysPresent: 0,
				DaysAbsent:  0,
				HoursWorked: 0,
			}
		}
		if !a.CheckIn.IsZero() && !a.CheckOut.IsZero() {
			hours := a.CheckOut.Sub(a.CheckIn).Hours()
			reportMap[a.WorkerID].HoursWorked += hours

			if hours > 0 {
				reportMap[a.WorkerID].DaysPresent++
			} else {
				reportMap[a.WorkerID].DaysAbsent++
			}
		} else {
			reportMap[a.WorkerID].DaysAbsent++
		}
	}

	reports := make([]report.WorkerAttendanceReport, 0, len(reportMap))
	for _, r := range reportMap {
		totalDays := r.DaysPresent + r.DaysAbsent
		if totalDays > 0 {
			r.AttendanceRate = float64(r.DaysPresent) / float64(totalDays) * 100
		}

		reports = append(reports, *r)
	}
	return reports, nil
}
