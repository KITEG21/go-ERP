package workers

import "user_api/internal/departments"

type Worker struct {
	ID           int                    `json:"id" gorm:"primaryKey"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	DepartmentId int                    `json:"department_id"`
	Department   departments.Department `json:"department" gorm:"foreignKey:DepartmentId;references:ID"`
	Salary       float32                `json:"salary"`
	HireDate     string                 `json:"hire_date"`
}
