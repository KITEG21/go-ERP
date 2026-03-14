package dtos

type CreateDepartmentDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
