package attendance

type AttendanceRepo interface {
	GetAllAttendance() ([]Attendance, error)
	FindPaginated(limit int, offset int) ([]Attendance, error)
	Count() (int64, error)
	CreateAttendance(attendance *Attendance) error
	UpdateAttendance(attendance *Attendance) error
	GetAttendanceById(id int) (Attendance, error)
	GetAttendancesByWorkerId(id int) ([]Attendance, error)
}
