package dtos

type UpdateAttendanceDto struct {
	ID       int     `json:"id" binding:"required"`
	CheckOut *string `json:"check_out,omitempty"`
}
