package payroll

import "github.com/rs/zerolog"

type PayrollService struct {
	repository PayrollRepo
	Logger     zerolog.Logger
}

func NewPayrollService(repository PayrollRepo, log zerolog.Logger) *PayrollService {
	return &PayrollService{repository: repository, Logger: log}
}

func (s *PayrollService) GetAllPayrolls() ([]Payroll, error) {
	return s.repository.GetAllPayrolls()
}

func (s *PayrollService) GetPayrollsPaginated(page int, pageSize int) ([]Payroll, int64, error) {
	count, err := s.repository.Count()
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	payrolls, err := s.repository.FindPaginated(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return payrolls, count, nil
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
