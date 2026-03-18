package dtos

type CreatePayrollDto struct {
	WorkerID   int     `json:"worker_id" validate:"required,gt=0"`
	Month      string  `json:"month" validate:"required"`
	BaseSalary float32 `json:"base_salary" validate:"required,gt=0"`
	Bonus      float32 `json:"bonus" validate:"gte=0"`
	Deductions float32 `json:"deductions"`
}
