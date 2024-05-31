package models

import "time"

//представление выражения(GET)
type ArithmeticRequest struct {
	Id int 				`json:"id"`
	Expression string   `json:"expression"`
}
//представление выражения(POST)
type Expression struct {
	Id int 			`json: 'id'`
	Status string   `json: 'status'`
	Result string   `json: 'result'`
}

//представление задачи для агента
type Task struct {
	Id int 						 `json: 'id'`
	Arg1 float64 				 `json: 'arg1'`
	Arg2 float64 				 `json: 'arg2'`
	Operation string 			 `json: 'operation'`
	Operation_time time.Duration `json: 'operation_time'`
}
//представление результата задачи
type ResultTask struct {
	Id int 	 		`json: 'id'`
	Result float64  `json: 'result'`
}

//представление постифксной формы выражения
type PostfixExpression struct {
	Id int
	Expression []string
}
