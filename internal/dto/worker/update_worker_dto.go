package dtos

type UpdateWorkerDto struct {
	Id    int    `json:"id" validate:"required,gt=0"`
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email"`
}
