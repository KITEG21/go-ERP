package dtos

type UpdateDepartmentDto struct {
	Id          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,notblank,min=3,max=100"`
	Description string `json:"description" validate:"max=255"`
}
