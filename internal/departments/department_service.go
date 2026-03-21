package departments

import "github.com/rs/zerolog" // Import zerolog

type DepartmentService struct {
	repo   DepartmentRepo
	Logger zerolog.Logger // Add Logger to the struct
}

func NewDepartmentService(repo DepartmentRepo, log zerolog.Logger) *DepartmentService { // Add log parameter
	return &DepartmentService{repo: repo, Logger: log}
}

func (s *DepartmentService) GetAllDepartments() ([]Department, error) {
	s.Logger.Info().Msg("Fetching all departments")
	return s.repo.GetAllDepartments()
}

func (s *DepartmentService) GetDepartmentsPaginated(page int, pageSize int) ([]Department, int64, error) {
	s.Logger.Info().Int("page", page).Int("pageSize", pageSize).Msg("Fetching paginated departments")
	count, err := s.repo.Count()
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to count departments")
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	departments, err := s.repo.FindPaginated(pageSize, offset)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to fetch paginated departments from repository")
		return nil, 0, err
	}
	return departments, count, nil
}

func (s *DepartmentService) CreateDepartment(department *Department) error {
	s.Logger.Info().Str("departmentName", department.Name).Msg("Creating new department")
	return s.repo.CreateDepartment(department)
}

func (s *DepartmentService) GetDepartmentById(id int) (Department, error) {
	s.Logger.Info().Int("departmentID", id).Msg("Fetching department by ID")
	return s.repo.GetDepartmentById(id)
}

func (s *DepartmentService) UpdateDepartment(department *Department) error {
	s.Logger.Info().Int("departmentID", int(department.ID)).Str("departmentName", department.Name).Msg("Updating department")
	return s.repo.UpdateDepartment(department)
}

func (s *DepartmentService) DeleteDepartment(id int) error {
	s.Logger.Info().Int("departmentID", id).Msg("Deleting department")
	return s.repo.DeleteDepartment(id)
}
