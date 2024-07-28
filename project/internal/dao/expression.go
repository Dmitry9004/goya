package dao

import (
	"goya/project/internal/models"
	"log"
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ExpressionDAO struct {
	ctx context.Context
	db *sql.DB
}

func NewExpressionDAO(db *sql.DB, ctx context.Context,) *ExpressionDAO {
	return &ExpressionDAO {
		ctx: ctx,
		db:db,
	}
}

func (dao *ExpressionDAO) CreateExpressionsTable() error {
	expressiosnTable := `
		CREATE TABLE IF NOT EXISTS expressions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			status TEXT,
			result TEXT,
			
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
	`
	
	if _, err := dao.db.ExecContext(dao.ctx, expressiosnTable); err != nil {
		return err
	}
	log.Println("ALL OK EXPRESSION")
	return nil
}

func (dao *ExpressionDAO)DeleteTable() {
	dao.db.ExecContext(dao.ctx, "DROP TABLE expressions")
}

func (dao *ExpressionDAO)SaveExpression(expression *models.Expression) (int, error) {
	query := `
		INSERT INTO expressions (user_id, status, result) VALUES ($1, $2, $3)
	`
	
	res, err := dao.db.ExecContext(dao.ctx, query, expression.UserId, expression.Status, expression.Result)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return int(id), nil
}

func (dao *ExpressionDAO)GetExpressionByIdAndUserId(id int, userId int) (models.Expression, error) {
	query := `
		SELECT * FROM expressions WHERE id = $1 AND user_id = $2
	`
	
	var expression models.Expression 
	err := dao.db.QueryRowContext(dao.ctx, query, id, userId).Scan(&expression.Id, &expression.UserId, &expression.Status, &expression.Result,)
	
	if err != nil {
		return expression, err
	}

	return expression, nil
}

func (dao *ExpressionDAO)GetExpressionById(id int) (models.Expression, error) {
	query := `
		SELECT * FROM expressions WHERE id = $1
	`
	
	var expression models.Expression 
	err := dao.db.QueryRowContext(dao.ctx, query, id).Scan(&expression.Id, &expression.UserId, &expression.Status, &expression.Result,)
	
	if err != nil {
		return expression, err
	}

	return expression, nil
}

func (dao *ExpressionDAO)GetAllExpressionByUserId(userId int) ([]models.Expression, error) {
	
	query := `
		SELECT * FROM expressions WHERE user_id = $1 
	`
	
	expressions := []models.Expression{}
	rows, err := dao.db.QueryContext(dao.ctx, query, userId)
	if err != nil {
		return expressions, err
	}
	for rows.Next() {
	defer rows.Close()
		var expression models.Expression
		err := rows.Scan(&expression.Id, &expression.UserId, &expression.Status, &expression.Result)
		
		if err != nil {
			return expressions, err
		}
		
		expressions = append(expressions, expression)
	}
	
	return expressions, nil
}

func (dao *ExpressionDAO)UpdateExpression(expression *models.Expression) error {
	query := `
		UPDATE expressions SET status = $1, result = $2 WHERE id = $3
	`
	
	if _, err := dao.db.ExecContext(dao.ctx, query, expression.Status, expression.Result, expression.Id); err != nil {
		return err
	}
	
	return nil
}