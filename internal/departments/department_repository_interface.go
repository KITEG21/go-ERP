package departments

type DepartmentRepo interface {
	GetAllDepartments() ([]Department, error)
	CreateDepartment(*Department) error
	GetDepartmentById(int) (Department, error)
	UpdateDepartment(*Department) error
	DeleteDepartment(int) error
}
