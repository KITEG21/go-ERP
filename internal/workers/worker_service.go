package workers

type WorkerService struct {
	repo WorkerRepository
}

func NewWorkerService(repo WorkerRepository) *WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) GetAllWorkers() ([]Worker, error) {
	return s.repo.GetAllWorkers()
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
