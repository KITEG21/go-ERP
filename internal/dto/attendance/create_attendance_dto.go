package dtos

import "time"

type CreateAttendanceDto struct {
	WorkerID int       `json:"worker_id" binding:"required"`
	CheckIn  time.Time `json:"check_in" binding:"required"`
	CheckOut time.Time `json:"check_out" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
}
