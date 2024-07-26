package services

import (
	"goya/project/internal/models"
	"goya/project/internal/dao"
	"goya/project/internal/services"
	
	"testing"
	"time"
	"database/sql"
	"context"
	"strings"
	"sync"
)

func TestToSimpleTask(t *testing.T) {
	postfixChan := make(chan *models.PostfixExpression)
	operationTime := map[string]time.Duration{}
	var mu *sync.Mutex
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "goya-test.db")
	if err != nil {
		t.Fatalf("error open database")
	}
	
	taskDAO := dao.NewTaskDAO(mu, db, ctx)
	expressionDAO := dao.NewExpressionDAO(db, ctx)
	
	go services.ToSimpleTask(postfixChan, operationTime, taskDAO, expressionDAO)
	
	expression := &models.Expression {
		UserId: 87,
		Result: "",
		Status: "process",
	}
	expressionDAO.CreateExpressionsTable()
	id, _ := expressionDAO.SaveExpression(expression)
	
	failExpression := &models.PostfixExpression {
		Id: id,
		Expression: strings.Split("THIS IS . FAIL STRING . EMPTY ", "."),
	}
	
	postfixChan <- failExpression
	
	time.Sleep(time.Second * 5)
	
	expressionEX, _ := expressionDAO.GetExpressionById(id)
	
	if expressionEX.Status != "not valid expression" {
		t.Fatalf("error from test to simple task")
	}
	
	expressionDAO.DeleteTable()
}