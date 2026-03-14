package attendance

import (
	"time"
	"user_api/internal/workers"
)

type Attendance struct {
	ID       int             `json:"id" gorm:"primaryKey"`
	WorkerID int             `json:"worker_id"`
	Worker   workers.Worker  `json:"worker" gorm:"foreignKey:WorkerId;references:ID"`
	CheckIn  time.Time       `json:"check_in" gorm:"type:time"`
	CheckOut time.Time       `json:"check_out" gorm:"type:time"`
	Date     time.Time       `json:"date" gorm:"type:date"`
	Status   AttendanceState `json:"status"`
}

type AttendanceState int

const (
	CheckedIn AttendanceState = iota
	CheckedOut
	Expired
)
