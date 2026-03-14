package workers

import (
	"user_api/internal/database"
)

type WorkerRepository struct{}

func (r *WorkerRepository) GetAllWorkers() ([]Worker, error) {
	var Workers []Worker
	error := database.DB.Find(&Workers).Error
	return Workers, error
}

func (r *WorkerRepository) FindPaginated(limit int, offset int) ([]Worker, error) {

	var workers []Worker

	err := database.DB.
		Limit(limit).
		Offset(offset).
		Find(&workers).Error

	return workers, err
}

func (r *WorkerRepository) CreateWorker(worker *Worker) error {
	return database.DB.Create(worker).Error
}

func (r *WorkerRepository) GetWorkerById(id int) (Worker, error) {
	var worker Worker
	error := database.DB.First(&worker, id).Error
	return worker, error
}

func (r *WorkerRepository) UpdateWorker(worker *Worker) error {
	return database.DB.Save(worker).Error
}

func (r *WorkerRepository) DeleteWorker(id int) error {
	return database.DB.Delete(&Worker{}, id).Error
}
