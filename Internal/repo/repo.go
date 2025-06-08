package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"micro-service/Internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	insertTaskQuery = `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`

	selectTaskByIDQuery = `SELECT 
    id,
    title, 
    description, 
    status
FROM tasks WHERE id = $1;`

	updateTaskQuery = `UPDATE tasks SET 
    title = $1,
    description = $2
WHERE id = $3
returning id;`

	deleteTaskQuery = `DELETE FROM tasks WHERE id = $1;`

	selectAllTasksQuery = `SELECT id, title, description, status FROM tasks;`
)

type repository struct {
	pool *pgxpool.Pool
}

type Repo interface {
	CreateTask(ctx context.Context, task Task) (int, error) // Создание задачи
	GetTask(ctx context.Context, id int) (Task, error)
	DeleteTask(ctx context.Context, id int) (int, error)
	PutTask(ctx context.Context, id int, task Task) (int, error)
	GetTasks(ctx context.Context) ([]Task, error)
}

func NewRepo(ctx context.Context, cfg config.PostgreSQL) (Repo, error) {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Parse PostgreSQL config")
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

/*// PatchTask - обновление задачи
func (r *repository) PatchTask(ctx context.Context, task Task) (int, error) {
	var id int
	fmt.Println("CTX: ", ctx, "\ntska: ", task)
	err := r.pool.QueryRow(ctx, patchTaskQuery, task.Title, task.Description, task.Id).Scan(&id)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, errors.Wrap(err, "failed to query patch task")
	}
	fmt.Println(id)
	return id, nil

}*/

// CreateTask - вставка новой задачи в таблицу tasks
func (r *repository) CreateTask(ctx context.Context, task Task) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert task")
	}
	return id, nil
}

// GetTasks - получение задачи
func (r *repository) GetTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, selectAllTasksQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query tasks")
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan task")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error during rows iteration")
	}

	return tasks, nil
}

// GetTask - получение задачи
func (r *repository) GetTask(ctx context.Context, id int) (Task, error) {
	var task Task
	err := r.pool.QueryRow(ctx, selectTaskByIDQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
	)
	if err != nil {
		return Task{}, errors.Wrap(err, "failed to take task")
	}
	return task, nil
}

// UpdateTask(put) - полное обновление
func (r *repository) PutTask(ctx context.Context, id int, task Task) (int, error) {
	err := r.pool.QueryRow(ctx, updateTaskQuery, task.Title, task.Description, id).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to update task")
	}
	return id, nil
}

// DeleteTask - удаление задачи
func (r *repository) DeleteTask(ctx context.Context, id int) (int, error) {
	result, err := r.pool.Exec(ctx, deleteTaskQuery, id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete task")
	}

	// Проверяем, сколько строк было удалено
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return 0, errors.New("task not found")
	}

	return id, nil
}
