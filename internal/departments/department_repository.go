package departments

import (
	"user_api/internal/database"
)

type DepartmentRepository struct{}

func (r *DepartmentRepository) GetAllDepartments() ([]Department, error) {
	var departments []Department
	error := database.DB.Find(&departments).Error
	return departments, error
}

func (r *DepartmentRepository) CreateDepartment(department *Department) error {
	return database.DB.Create(department).Error
}

func (r *DepartmentRepository) GetDepartmentById(id int) (Department, error) {
	var department Department
	error := database.DB.First(&department, id).Error
	return department, error
}

func (r *DepartmentRepository) UpdateDepartment(department *Department) error {
	return database.DB.Save(department).Error
}

func (r *DepartmentRepository) DeleteDepartment(id int) error {
	return database.DB.Delete(&Department{}, id).Error
}
