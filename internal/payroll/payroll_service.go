package payroll

type PayrollService struct {
	repository PayrollRepository
}

func NewPayrollService(repository PayrollRepository) *PayrollService {
	return &PayrollService{repository: repository}
}

func (s *PayrollService) GetAllPayrolls() ([]Payroll, error) {
	return s.repository.GetAllPayrolls()
}

func (s *PayrollService) CreatePayroll(payroll *Payroll) error {
	return s.repository.CreatePayroll(payroll)
}

func (s *PayrollService) GetPayrollById(id int) (Payroll, error) {
	return s.repository.GetPayrollById(id)
}

func (s *PayrollService) UpdatePayroll(payroll *Payroll) error {
	return s.repository.UpdatePayroll(payroll)
}

func (s *PayrollService) GetPayrollsByWorkerId(id int) ([]Payroll, error) {
	return s.repository.GetPayrollByWorkerId(id)
}
