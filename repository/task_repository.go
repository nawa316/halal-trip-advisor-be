package repository

import (
	"context"
	"database/sql"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	_, err := tr.db.ExecContext(c,
		`INSERT INTO tasks (id, title, user_id) VALUES ($1, $2, $3)`,
		task.ID,
		task.Title,
		task.UserID,
	)

	return err
}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	rows, err := tr.db.QueryContext(c, `SELECT id, title, user_id FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.UserID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}