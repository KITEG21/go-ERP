package payroll

import (
	"user_api/internal/database"
)

type PayrollRepository struct{}

func (r *PayrollRepository) GetAllPayrolls() ([]Payroll, error) {
	var payrolls []Payroll
	error := database.DB.Preload("Worker").Find(&payrolls).Error
	return payrolls, error
}

func (r *PayrollRepository) CreatePayroll(payroll *Payroll) error {
	return database.DB.Create(payroll).Error
}

func (r *PayrollRepository) GetPayrollById(id int) (Payroll, error) {
	var payroll Payroll
	error := database.DB.Preload("Worker").First(&payroll, id).Error
	return payroll, error
}

func (r *PayrollRepository) UpdatePayroll(payroll *Payroll) error {
	return database.DB.Save(payroll).Error
}

func (r *PayrollRepository) DeletePayroll(id int) error {
	return database.DB.Delete(&Payroll{}, id).Error
}

func (r *PayrollRepository) GetPayrollByWorkerId(id int) ([]Payroll, error) {
	var payrolls []Payroll
	error := database.DB.Where("worker_id = ?", id).Preload("Worker").Find(&payrolls).Error
	return payrolls, error
}

func (r *PayrollRepository) FindPaginated(limit int, offset int) ([]Payroll, error) {
	var payrolls []Payroll
	err := database.DB.Preload("Worker").
		Limit(limit).
		Offset(offset).
		Find(&payrolls).Error
	return payrolls, err
}

func (r *PayrollRepository) Count() (int64, error) {
	var count int64
	err := database.DB.Model(&Payroll{}).Count(&count).Error
	return count, err
}
