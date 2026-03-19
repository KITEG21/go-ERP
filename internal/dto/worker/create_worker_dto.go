package dtos

type CreateWorkerDto struct {
	Name         string  `json:"name" validate:"required,min=3,max=100,notblank"`
	Email        string  `json:"email" validate:"required,email"`
	DepartmentId *int    `json:"department_id,omitempty" validate:"omitempty,gt=0"`
	HireDate     *string `json:"hire_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}
