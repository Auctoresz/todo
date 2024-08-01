package dao

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("dao: no matching record found")
)

type Task struct {
	Id           int64
	IdCustomer   int64
	Description  string
	DueDate      time.Time
	IsCompleted  bool
	CreationDate time.Time
}
