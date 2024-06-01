package main
import (
	"goya/project/internal/config"
	"goya/project/internal/models"
	"goya/project/internal/services" 
	
    "log"
	"net/http"
	"encoding/json"
	"strconv"
	
	"github.com/go-chi/chi/v5"
	"github.com/recoilme/pudge"
)

func main() {
	countExec, operationsTime := config.SetConfig()

	//каналы для представлений 
	chanPostfixExpression := make(chan *models.PostfixExpression)
	chanResults := make(chan *models.ResultTask)
	chanTasks   := make(chan models.Task)
	
    //встроенная база данных для складирования представлений
    db, _                := pudge.Open("../expressions", nil)
    dbTasksRaw, _        := pudge.Open("../tasks/raw", nil)
    dbTasksInProccess, _ := pudge.Open("../tasks/proccess", nil)
	dbTasksDone, _       := pudge.Open("../tasks/done", nil)
	
	defer db.DeleteFile()
	defer dbTasksRaw.DeleteFile()
	defer dbTasksInProccess.DeleteFile()
	defer dbTasksDone.DeleteFile()
	
	//запуск агентов в отдельных выражений до переменной COMPUTING_POWER 
	for i := 0; i < countExec; i++ {
		go services.ExecTask(chanTasks, chanResults)
	}
	
	//запуск метода преобразователя в отедльной горутине
	go services.ToSimpleTask(chanPostfixExpression, operationsTime)
	
	//запуск метода отправления результатов задач серверу в отедльной горутине
	go services.WaiterResults(chanResults)

	//запуск метода приема задач с сервера
	go services.TakerTask(chanTasks)

	r := chi.NewRouter()
	//endpoint для приема выражения с клиента
	r.Post("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		request := models.ArithmeticRequest{}
		
		err := json.NewDecoder(r.Body).Decode(&request);
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		
		if (len(request.Expression) < 3) {
		    log.Println("error: /api/v1/calculate, POST, " + request.Expression)
			w.WriteHeader(422)
			return
		}
		
		postfixExppression, errorParse := services.ToPostfix(request.Expression)
		
		w.Header().Set("Content-Type", "applciation/json")
		if errorParse != nil {
			response := map[string]string {
				"error": "not valide expression",
			}
			json.NewEncoder(w).Encode(response)
		    log.Println("error: /api/v1/calculate, POST, " + request.Expression)
		    w.WriteHeader(422)
		    return
		}
		
		id := services.GenerateId()
		
		postfix := &models.PostfixExpression {
			Id: id,
			Expression: postfixExppression,
		}
		
	    chanPostfixExpression <- postfix
	    
		expression := &models.Expression {
			Id: id,
			Status: "doing",
			Result: "none",
		}
		
		db.Set(id, expression)
		response := map[string]int {
			"id":id,
		}

		json.NewEncoder(w).Encode(response)
	    w.WriteHeader(201)
	})	
	
	////endpoint для отправки всех выражений клиенту
	r.Get("/api/v1/expressions", func(w http.ResponseWriter, r *http.Request) {
		expressions := []models.Expression{}
		keys, _ := db.Keys(nil, 0, 0, false)
		
		if len(keys) == 0 {
		    w.WriteHeader(404)
		    log.Println("error: /api/v1/expressions, GET ")
		    return
		}
		
		for _, v := range keys {
		    var expression models.Expression
            db.Get(v, &expression)
            
            expressions = append(expressions, expression)
		}
		
		w.Header().Set("Content-Type", "applciation/json")
		json.NewEncoder(w).Encode(expressions)
		w.WriteHeader(200)
	})
	
	//endpoint для приема выражжения с определенным ID
	r.Get("/api/v1/expressions/{id}", func(w http.ResponseWriter, r *http.Request) {
        var resultExpression models.Expression
		idParameter := chi.URLParam(r, "id")
		keys, _ := db.Keys(nil, 0, 0, false)
		
		for _, v := range keys {
		    var expression models.Expression
            db.Get(v, &expression)
            
            if strconv.Itoa(expression.Id) == idParameter {
                resultExpression = expression
                break
            }
		}
		
		if resultExpression.Id == 0 {
		    w.WriteHeader(404)
		    log.Println("error: /api/v1/expressions, GET, id: " + idParameter)
		    return
		}
		
		w.Header().Set("Content-Type", "applciation/json")
		json.NewEncoder(w).Encode(resultExpression)
		w.WriteHeader(200)
	})
	
	//endpoint для отправки "сырой" задачи агенту
	r.Get("/internal/task", func(w http.ResponseWriter, r *http.Request) {
	    var task models.Task
	    keys, err := dbTasksRaw.Keys(nil, 0, 0, false)
	    if err != nil {
			log.Println(err)
		}
		
		if len(keys) == 0 {
	       w.WriteHeader(404)
	       log.Println("info: empty raw tasks")
	       return
	    }
		
		key := string(keys[0])
		
	    dbTasksRaw.Get(key, &task)
	    dbTasksRaw.Delete(key)
	    
		dbTasksInProccess.Set(key, task)
		
		w.Header().Set("Content-Type", "applciation/json")
	    json.NewEncoder(w).Encode(&task)
	})
	
	//endpoint для преима готовой задачи с агента
	r.Post("/internal/task", func(w http.ResponseWriter, r *http.Request) {
	    var task models.ResultTask 
	    var guessTask models.Task
	    json.NewDecoder(r.Body).Decode(&task)
		
	    dbTasksInProccess.Get(task.Id, &guessTask)
	    if task.Id != guessTask.Id {
	        w.WriteHeader(404)
	        log.Println("error: /internal/task, POST, description: not found task")
	    }
		
	    dbTasksDone.Set(task.Id, task)
	})
	

	http.ListenAndServe(":8080", r)
}
