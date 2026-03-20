package departments

type DepartmentRepo interface {
	FindPaginated(limit int, offset int) ([]Department, error)
	GetAllDepartments() ([]Department, error)
	CreateDepartment(*Department) error
	GetDepartmentById(int) (Department, error)
	UpdateDepartment(*Department) error
	DeleteDepartment(int) error
	Count() (int64, error)
}
