package dtos

type ResponsePayrollDto struct {
	Id         int     `json:"id"`
	WorkerID   int     `json:"worker_id" binding:"required"`
	Month      string  `json:"month" binding:"required"`
	BaseSalary float32 `json:"base_salary"`
	Bonus      float32 `json:"bonus"`
	Deductions float32 `json:"deductions"`
	NetSalary  float32 `json:"net_salary"`
}
