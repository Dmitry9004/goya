package main

import "time"
import "net/http"
import "encoding/json"
import "log"
import "errors"
//import "sync"
import "os"
import "strconv"

import "github.com/go-chi/chi/v5"

//variable
//import "github.com/joho/godotenv"

//представление задачи для агента
type Task struct {
	Id int 						 "json: 'id'"
	Arg1 float64 				 "json: 'arg1'"
	Arg2 float64 				 "json: 'arg2'"
	Operation string 			 "json: 'operation'"
	Operation_time time.Duration "json: 'operation_time'"
}

type ResultTask struct {
	Id int 	 		"json: 'id'"
	Result float64  "json: 'result'"
}

func main() {
	//count gorutines from .env
	count, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))

	//channels for get Task and post results 
	chanResults := make(chan *ResultTask)
	chanTasks := make(chan Task)

	//run goroutines for execution task 
	for i := 0; i < count; i++ {
		go execTask(chanTasks, chanResults)
	}

	//main agent router
	r := chi.NewRouter()

	//wait results from chanResults
	go waiterResults(chanResults, r)


	log.Println("-- -- RUN IT -- --")

	//pool workers
	r.Post("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(404)
			log.Print("task: POST -- 404")
		}
		//task for our execTask chan
		var task Task

		//
		errorDecode := json.NewDecoder(r.Body).Decode(&task);
		if errorDecode != nil {
			w.WriteHeader(500)
			log.Print("task: POST -- 500")
		}
		
		chanTasks <- task

		//need get error if operation not in our container operations
	})

	http.ListenAndServe(":8080", r)
}

func execTask(chanTasks chan Task, chanResults chan *ResultTask) {
	for {
		select {
			//if chan has task
			case task := <- chanTasks:
				//exec it
				result, err := calculate(task) 
				if err != nil {
					//what need do if errore?????????????
					log.Print(err)					
				}

				//result for send to server
				responseResult := &ResultTask {
					Id: task.Id,
					Result: result,
				}

				//need send to server, send to chan
				chanResults <- responseResult
		}
	}
}

func waiterResults(chanResults chan *ResultTask, r *chi.Mux) {
	for {
		select {
			case result := <- chanResults:
				r.Post("/internal/task", func(w http.ResponseWriter, r *http.Request) {
					//post result to endpoint

					response := map[string]interface{} {
						"data": map[string]interface{} {
							"id": result.Id,
							"result": result.Result,
						},
					}

		//response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
				})
			//case timer 
		}
	}
}

//function execute task's operation
func calculate(task Task) (float64, error) {	
	
	result := 0.0
	switch(task.Operation) { 
		case "+": 
			result = task.Arg1 + task.Arg2
		case "-":
			result = task.Arg1 - task.Arg2
		case "*":
			result = task.Arg1 * task.Arg2
		case "/":
			result = task.Arg1 / task.Arg2
		default:
			log.Print("not find operation: " + string(task.Operation) + "id: " + string(task.Id))
			return 0, errors.New("not find operation") 
	}

	return result, nil
}