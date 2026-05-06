package domain

import (
	"context"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     string `json:"id,omitempty"`
	Title  string `form:"title" binding:"required" json:"title"`
	UserID string `json:"-"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
