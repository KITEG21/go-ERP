package attendance

import (
	"user_api/internal/database"
)

type AttendanceRepository struct{}

func (r *AttendanceRepository) GetAllAttendance() ([]Attendance, error) {
	var attendances []Attendance
	err := database.DB.Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindPaginated(limit int, offset int) ([]Attendance, error) {
	var attendances []Attendance
	err := database.DB.
		Limit(limit).
		Offset(offset).
		Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) Count() (int64, error) {
	var count int64
	err := database.DB.Model(&Attendance{}).Count(&count).Error
	return count, err
}

func (r *AttendanceRepository) CreateAttendance(attendance *Attendance) error {
	return database.DB.Create(attendance).Error
}

func (r *AttendanceRepository) UpdateAttendance(attendance *Attendance) error {
	return database.DB.Save(attendance).Error
}

func (r *AttendanceRepository) GetAttendanceById(id int) (Attendance, error) {
	var attendance Attendance
	err := database.DB.First(&attendance, id).Error
	return attendance, err
}

func (r *AttendanceRepository) GetAttendancesByWorkerId(id int) ([]Attendance, error) {
	var attendances []Attendance
	err := database.DB.Where("worker_id =?", id).Find(&attendances).Error
	return attendances, err
}
