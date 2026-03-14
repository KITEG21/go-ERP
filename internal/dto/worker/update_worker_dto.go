package dtos

type UpdateWorkerDto struct {
	Id    int    `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
