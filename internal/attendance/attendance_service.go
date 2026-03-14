package attendance

type AttendanceService struct {
	repo *AttendanceRepository
}

func NewAttendanceService(repo *AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) GetAllAttendance() ([]Attendance, error) {
	return s.repo.GetAllAttendance()
}

func (s *AttendanceService) CreateAttendance(attendance *Attendance) error {
	return s.repo.CreateAttendance(attendance)
}

func (s *AttendanceService) UpdateAttendance(attendance *Attendance) error {
	return s.repo.UpdateAttendance(attendance)
}

func (s *AttendanceService) GetAttendanceById(id int) (Attendance, error) {
	return s.repo.GetAttendanceById(id)
}

func (s *AttendanceService) GetAttendancesByWorkerId(id int) ([]Attendance, error) {
	return s.repo.GetAttendancesByWorkerId(id)
}
