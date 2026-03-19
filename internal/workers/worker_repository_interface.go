package workers

type WorkerRepo interface {
	GetAllWorkers() ([]Worker, error)
	FindPaginated(limit int, offset int) ([]Worker, error)
	Count() (int64, error)
	CreateWorker(*Worker) error
	GetWorkerById(id int) (Worker, error)
	UpdateWorker(*Worker) error
	DeleteWorker(id int) error
}
