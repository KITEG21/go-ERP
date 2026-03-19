package departments

type DepartmentService struct {
	repo DepartmentRepository
}

func NewDepartmentService(repo DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) GetAllDepartments() ([]Department, error) {
	return s.repo.GetAllDepartments()
}

func (s *DepartmentService) GetDepartmentsPaginated(page int, pageSize int) ([]Department, int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	departments, err := s.repo.FindPaginated(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return departments, count, nil
}

func (s *DepartmentService) CreateDepartment(department *Department) error {
	return s.repo.CreateDepartment(department)
}

func (s *DepartmentService) GetDepartmentById(id int) (Department, error) {
	return s.repo.GetDepartmentById(id)
}

func (s *DepartmentService) UpdateDepartment(department *Department) error {
	return s.repo.UpdateDepartment(department)
}

func (s *DepartmentService) DeleteDepartment(id int) error {
	return s.repo.DeleteDepartment(id)
}
