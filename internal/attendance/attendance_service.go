package attendance

import "github.com/rs/zerolog" // Import zerolog

type AttendanceService struct {
	repo   *AttendanceRepository
	Logger zerolog.Logger // Add Logger to the struct
}

func NewAttendanceService(repo *AttendanceRepository, log zerolog.Logger) *AttendanceService { // Add log parameter
	return &AttendanceService{repo: repo, Logger: log}
}

func (s *AttendanceService) GetAllAttendance() ([]Attendance, error) {
	s.Logger.Info().Msg("Fetching all attendance records") // Example log
	return s.repo.GetAllAttendance()
}

func (s *AttendanceService) GetAttendancesPaginated(page int, pageSize int) ([]Attendance, int64, error) {
	s.Logger.Info().Int("page", page).Int("pageSize", pageSize).Msg("Fetching paginated attendance records") // Example log
	count, err := s.repo.Count()
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to count attendance records") // Example error log
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	attendances, err := s.repo.FindPaginated(pageSize, offset)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to fetch paginated attendance records from repository") // Example error log
		return nil, 0, err
	}
	return attendances, count, nil
}

func (s *AttendanceService) CreateAttendance(attendance *Attendance) error {
	s.Logger.Info().Msg("Creating new attendance record") // Example log
	return s.repo.CreateAttendance(attendance)
}

func (s *AttendanceService) UpdateAttendance(attendance *Attendance) error {
	s.Logger.Info().Int("id", int(attendance.ID)).Msg("Updating attendance record") // Example log
	return s.repo.UpdateAttendance(attendance)
}

func (s *AttendanceService) GetAttendanceById(id int) (Attendance, error) {
	s.Logger.Info().Int("id", id).Msg("Fetching attendance record by ID") // Example log
	return s.repo.GetAttendanceById(id)
}

func (s *AttendanceService) GetAttendancesByWorkerId(id int) ([]Attendance, error) {
	s.Logger.Info().Int("workerID", id).Msg("Fetching attendance records by worker ID") // Example log
	return s.repo.GetAttendancesByWorkerId(id)
}
