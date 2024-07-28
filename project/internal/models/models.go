package models

import "time"

type User struct {
	Id int 			`json: "id"`
	Username string `json: "username"`
	Password string `json: "password"`
}

//представление выражения(GET)
type ArithmeticRequest struct {
	Expression string   `json:"expression"`
}
//представление выражения(POST)
type Expression struct {
	Id int 				`json: 'id'`
	UserId int 			`json: 'user_id'`
	Status string   	`json: 'status'`
	Result string   	`json: 'result'`
}

//представление задачи для агента
type Task struct {
	Id int 						 `json: 'id'`
	ExpressionId int 			 `json: 'expression_id'`
	Arg1 float64 				 `json: 'arg1'`
	Arg2 float64 				 `json: 'arg2'`
	Result float64					 `json: 'result'`
	Operation string 			 `json: 'operation'`
	OperationTime time.Duration `json: 'operation_time'`
	Status string 				 `json: 'status'`
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


func NewStack[T any]() *Stack[T] {
    return &Stack[T]{
        arr:[]T{},
    }
}

type Stack[T any] struct {
    arr []T
    index int
}

func (st *Stack[T]) Len() int {
    return st.index
}

func (st *Stack[T]) Push(str T) {
    st.arr = append(st.arr, str)
    st.index++
}

func (st *Stack[T]) Pop() T {
    val := st.arr[st.index-1]
    if st.index == 0 {
        st.arr = []T{}
        return val
    }
    
    st.index--
    st.arr = st.arr[:st.index]
    return val
}

func (st *Stack[T]) Peek() T {
    if st.index == 0 {
        var val T
        return val
    }
    return st.arr[st.index-1]
}