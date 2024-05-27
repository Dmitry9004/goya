package internal

import "encoding/json"
import "github.com/go-chi/chi/v5"
import "github.com/recoilme/pudge"
import "time"

import ""

type Expression struct {
	Id int 			  "json: 'id'"
	Expression string "json: 'expression'"
}

func main() {
	r := chi.NewRouter()

	pudge.Set("../expressions", "Test", "Hey")

	r.Post("/api/v1/calculate", func(w http.WriterResponse, r *http.Reqeust) {
		var expression Expression
		json.NewDecoder(r.Body).Decode(&expression)

		postfixExp := toPostfix(expression.Expression)
		toSimpleTask(postfixExp)
	})

	r.Get("/api/v1/expressions", func(w http.WriterResponse, r *http.Reqeust) {

	})

	r.Get("/api/v1/expressions{id}", func(w http.WriterResponse, r *http.Reqeust) {

	})	


	//

	r.Get("internal/task", func(w http.WriterResponse, r *http.Reqeust) {

	})

	r.Post("internal/task", func(w http.WriterResponse, r *http.Reqeust) {

	})

	http.ListenAndServe(":8080", r)	
}


func toPostfix(str string) string {
	return str
}

func toSimpleTask(str string) {
	idCount := 1
	st := NewStack()

	for {}
}

type Task struct {
	Id int 				"json: 'id'"
	Arg1 string 		"json: 'arg1'"
	Arg2 string 		"json: 'arg2'"
	Operation string 	"json: 'operation'"
	Operation_time time "json: 'operation_time'"
}