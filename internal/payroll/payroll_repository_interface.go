package payroll

type PayrollRepo interface {
	GetAllPayrolls() ([]Payroll, error)
	CreatePayroll(*Payroll) error
	GetPayrollById(int) (Payroll, error)
	UpdatePayroll(*Payroll) error
	DeletePayroll(int) error
	GetPayrollByWorkerId(int) ([]Payroll, error)
	FindPaginated(limit int, offset int) ([]Payroll, error)
	Count() (int64, error)
}
