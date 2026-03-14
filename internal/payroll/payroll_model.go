package payroll

import "user_api/internal/workers"

type Payroll struct {
	Id         int            `json:"id" gorm:"primaryKey"`
	WorkerId   int            `json:"worker_id"`
	Worker     workers.Worker `json:"worker" gorm:"foreignKey:WorkerId;references:ID"`
	Month      string         `json:"month"`
	BaseSalary float32        `json:"base_salary"`
	Bonus      float32        `json:"bonus"`
	Deductions float32        `json:"deductions"`
	NetSalary  float32        `json:"net_salary"`
	Status     PayrollStatus  `json:"status"`
}

type PayrollStatus int

const (
	Pending PayrollStatus = iota
	Processed
	Failed
)
