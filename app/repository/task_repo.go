package repository

import (
	"context"
	"database/sql"
	"errors"
	"todo-echo/app/domain/dao"
)

type TaskRepository interface {
	GetAll(ctx context.Context, task dao.Task) ([]dao.Task, error)
	Get(ctx context.Context, task dao.Task) (dao.Task, error)
	Insert(ctx context.Context, task dao.Task) (dao.Task, error)
	Delete(ctx context.Context, task dao.Task) error
	Completed(ctx context.Context, task dao.Task) error
	Edit(ctx context.Context, task dao.Task) error
}

type TaskRepositoryImpl struct {
	Db *sql.DB
}

// NewTaskRepositoryImpl is a Constructor
func NewTaskRepositoryImpl(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{
		Db: db,
	}
}

func (t TaskRepositoryImpl) GetAll(ctx context.Context, task dao.Task) ([]dao.Task, error) {
	var tasks []dao.Task

	script := "SELECT id, id_customer, description, due_date, is_completed, creation_date FROM task WHERE id_customer=?"
	rows, err := t.Db.QueryContext(ctx, script, task.IdCustomer)
	// error handling
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&task.Id, &task.IdCustomer, &task.Description, &task.DueDate, &task.IsCompleted, &task.CreationDate)
		// error handling
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t TaskRepositoryImpl) Get(ctx context.Context, task dao.Task) (dao.Task, error) {
	script := "SELECT id, id_customer, description, due_date, is_completed, creation_date FROM task WHERE id = ?"
	row := t.Db.QueryRowContext(ctx, script, task.Id)

	err := row.Scan(&task.Id, &task.IdCustomer, &task.Description, &task.DueDate, &task.IsCompleted, &task.CreationDate)
	// error handling
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, dao.ErrNoRecord
		} else {
			return task, err
		}
	}

	return task, nil
}

func (t TaskRepositoryImpl) Insert(ctx context.Context, task dao.Task) (dao.Task, error) {
	script := "INSERT INTO task(id_customer, description, due_date) VALUES (?, ?, ?)"
	res, err := t.Db.ExecContext(ctx, script, task.IdCustomer, task.Description, task.DueDate)
	// error handling
	if err != nil {
		return task, err
	}

	id, err := res.LastInsertId()
	// error handling
	if err != nil {
		return task, err
	}

	task.Id = id
	return task, nil
}

func (t TaskRepositoryImpl) Delete(ctx context.Context, task dao.Task) error {
	script := "DELETE FROM task WHERE id=? AND id_customer=?"
	_, err := t.Db.ExecContext(ctx, script, task.Id, task.IdCustomer)
	if err != nil {
		return err
	}
	return nil
}

func (t TaskRepositoryImpl) Completed(ctx context.Context, task dao.Task) error {
	script := "UPDATE task SET is_completed=? WHERE id=? AND id_customer=?"
	_, err := t.Db.ExecContext(ctx, script, task.IsCompleted, task.Id, task.IdCustomer)
	if err != nil {
		return err
	}

	return nil
}

func (t TaskRepositoryImpl) Edit(ctx context.Context, task dao.Task) error {
	script := "UPDATE task SET description=?, due_date=? WHERE id=? AND id_customer=?"
	_, err := t.Db.ExecContext(ctx, script, task.Description, task.DueDate, task.Id, task.IdCustomer)
	if err != nil {
		return err
	}

	return nil
}
