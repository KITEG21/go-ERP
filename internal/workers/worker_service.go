package workers

import "github.com/rs/zerolog"

type WorkerService struct {
	repo   WorkerRepo
	Logger zerolog.Logger
}

func NewWorkerService(repo WorkerRepo, log zerolog.Logger) *WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) GetAllWorkers() ([]Worker, error) {
	return s.repo.GetAllWorkers()
}

func (s *WorkerService) GetWorkersPaginated(page int, pageSize int) ([]Worker, int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	workers, err := s.repo.FindPaginated(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return workers, count, nil
}

func (s *WorkerService) CreateWorker(worker *Worker) error {
	return s.repo.CreateWorker(worker)
}
func (s *WorkerService) GetWorkerById(id int) (Worker, error) {
	return s.repo.GetWorkerById(id)
}

func (s *WorkerService) UpdateWorker(worker *Worker) error {
	return s.repo.UpdateWorker(worker)
}

func (s *WorkerService) DeleteWorker(id int) error {
	return s.repo.DeleteWorker(id)
}
