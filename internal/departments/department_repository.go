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

func (r *DepartmentRepository) FindPaginated(limit int, offset int) ([]Department, error) {
	var departments []Department
	err := database.DB.
		Limit(limit).
		Offset(offset).
		Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) Count() (int64, error) {
	var count int64
	err := database.DB.Model(&Department{}).Count(&count).Error
	return count, err
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
