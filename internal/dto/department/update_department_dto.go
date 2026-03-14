package dtos

type UpdateDepartmentDto struct {
	Id          int    `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
