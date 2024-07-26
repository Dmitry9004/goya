package dao

import (
	"goya/project/internal/models"

	"sync"
	"log"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


type TaskDAO struct {
	mu *sync.Mutex
	db *sql.DB
	ctx context.Context
}

func NewTaskDAO(mu *sync.Mutex, db *sql.DB, ctx context.Context) *TaskDAO {
	return &TaskDAO{
		mu:mu,
		db: db,
		ctx:ctx,
	}
} 

func (dao *TaskDAO)CreateTasksTable() error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	
	tasksTable := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER not null,
			expression_id INTEGER,
			arg1 TEXT,
			arg2 TEXT,
			result TEXT,
			operation TEXT,
			operation_time TEXT,
			status TEXT,
			
			FOREIGN KEY (expression_id) REFERENCES expressions (id)
		);
	`
	
	if _, err := dao.db.ExecContext(dao.ctx, tasksTable); err != nil {
		return err
	}
	log.Println("ALL OK TASK")
	return nil
}

func (dao *ExpressionDAO)DeletTable() {
	dao.db.ExecContext(dao.ctx, "DROP TABLE expressions")
}

func (dao *TaskDAO)SaveTask(task *models.Task) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	
	query := `
		INSERT INTO tasks (id, expression_id, status, arg1, arg2, result, operation, operation_time, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if _, err := dao.db.ExecContext(dao.ctx, query, task.Id, task.ExpressionId, task.Status, task.Arg1, task.Arg2, task.Result, task.Operation, task.OperationTime, task.Status); err != nil {
		//log.Println(err)
		return err
	}
	
	return nil
}

func (dao *TaskDAO)UpdateTask(task *models.Task) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	
	query := `
		UPDATE tasks SET status = $1, result = $2 WHERE id = $3
	`
	
	if _, err := dao.db.ExecContext(dao.ctx, query, task.Status, task.Result, task.Id); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (dao *TaskDAO)GetTask(id int) (models.Task, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	
	query := `
		SELECT * from tasks WHERE id = $1
	`
	var res models.Task
	err := dao.db.QueryRowContext(dao.ctx, query, id).Scan(&res.Id, &res.ExpressionId, &res.Arg1, &res.Arg2, &res.Result, &res.Operation, &res.OperationTime, &res.Status)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (dao *TaskDAO)GetRawTasks(status string) ([]models.Task, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	query := `
		SELECT * FROM tasks WHERE status = $1
	`
	
	resultTasks := []models.Task{}
	rows, err := dao.db.QueryContext(dao.ctx, query, status)
	
	if err != nil {
		return resultTasks, err
	}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.ExpressionId, &task.Arg1, &task.Arg2, &task.Result, &task.Operation, &task.OperationTime, &task.Status)
		if err != nil {
			//log.Println(err)
			return resultTasks, err
		}
		resultTasks = append(resultTasks, task)
	}
	
	return resultTasks, err
}

