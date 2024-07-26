package tests

import (
	"goya/project/internal/dao"
	"goya/project/internal/models"

	"log"
	"testing"
	"context"
	"database/sql"
	
	_ "github.com/mattn/go-sqlite3"
)

func openDatabase(nameDB string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		return db, err
		log.Fatalf("error open sqlite3")
	}
	return db, nil
}

func createTable(db *sql.DB) *expressionDAO {
	expressionDAO := dao.NewExpressionDAO(db, context.TODO())
	err = expressionDAO.CreateExpressionsTable();
	return expressionDAO
}

func TestMain(m *testing.M) {
	db, err := openDatabase("goya-test")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	dao := createTable(db)
	os.Args = []interface{}{dao}
	
	m.Run()
	
	dao.DeleteTable()
	//closeDatabase()
}

func TestSaveExpression(t *testing.T) {
	expressionDAO := os.Args[0].(*dao.ExpressionDAO)
	expression := models.Expression{
		Id: 532,
		UserId: 22,
		Result: "98345",
		Status: "done",
	}
	
	_, err = expressionDAO.SaveExpression(&expression)
	if err != nil {
		t.Fatalf("err")
	}
}

func TestGetExpressionByIdAndUserId(t *testing.T) {
	expressionDAO := os.Args[0].(*models.ExpressionDAO)
	expression := models.Expression {
		UserId: 754,
		Result: "",
		Status: "process",
	}
	
	expressionDAO.CreateExpressionsTable()
	
	id, err := expressionDAO.SaveExpression(&expression)
	
	expressionEx, err := expressionDAO.GetExpressionByIdAndUserId(id, 754)
	if err != nil {
		t.Fatalf("error from get expressions by id and user id")
	}
	if expressionEx.UserId != expression.UserId {
		t.Fatalf("error from from expression by id and user id")
	}
}

func TestGetAllExpressionByUserId(t *testing.T) {
	expressionDAO := os.Args[0].(*models.ExpressionDAO)
	expression := models.Expression{
		UserId: 754,
		Result: "",
		Status: "process",
	}
	
	expressionDAO.CreateExpressionsTable()
	
	_, err = expressionDAO.SaveExpression(&expression)
	expressions, err := expressionDAO.GetAllExpressionByUserId(754)
	if err != nil {
		t.Fatalf("error from get all expressions ")
	}
	
	if len(expressions) == 0 {
		t.Fatalf("error from get all")
	}
}