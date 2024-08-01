package dto

import "time"

type TaskCreateRequest struct {
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type TaskCompletedRequest struct {
	IsCompleted bool `validate:"required,boolean" json:"is_completed"`
}

type TaskGetRequest struct {
	Id int64 `validate:"required" uri:"taskId"`
}
