package main
import (
	"goya/project/internal/config"
	"goya/project/internal/models"
	"goya/project/internal/services" 
	"goya/project/internal/dao"
	"goya/project/internal/rpc"
	"goya/project/internal/auth"
	pb "goya/project/proto"
	
    "log"
	"net/http"
	"encoding/json"
	"strconv"
	"database/sql"
	"context"
	"fmt"
	"sync"
	"os"
	"net"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	countExec, operationsTime := config.SetConfig()
	
	//каналы для представлений 
	chanPostfixExpression := make(chan *models.PostfixExpression)
	chanResults := make(chan *models.ResultTask)
	chanTasks   := make(chan models.Task)
	
    //встроенная база данных для складирования представлений
    //db, _                := pudge.Open("../expressions", nil)
    //dbTasksRaw, _        := pudge.Open("../tasks/raw", nil)
    //dbTasksInProccess, _ := pudge.Open("../tasks/proccess", nil)
	//dbTasksDone, _       := pudge.Open("../tasks/done", nil)
	
	
	//sqlite3 
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "goya.db")
	if err != nil {
		log.Println("error open database")
	}
	defer db.Close()
	
	var mu sync.Mutex
	userDAO := dao.NewUserDAO(ctx, db)
	expressionDAO := dao.NewExpressionDAO(db, ctx)
	
	if err = userDAO.CreateUsersTable(); err != nil {
		log.Println(err)
		//exit
	}
	
	if err = expressionDAO.CreateExpressionsTable(); err != nil {
		log.Println(err)
	}
	
	taskDAO := dao.NewTaskDAO(&mu, db, ctx)
	
	if err = taskDAO.CreateTasksTable(); err != nil {
		log.Println(err)
	}	
	
	//запуск grpc сервера
	host := "localhost"
	port := "8000"
		
	addr := fmt.Sprintf("%s:%s", host, port)
		
	go func() {
		lis, _ := net.Listen("tcp", addr)
		
		rpcServer := grpc.NewServer()
		taskServer := rpc.NewServer(taskDAO)
		
		pb.RegisterTaskServiceServer(rpcServer, taskServer)
		
		if err := rpcServer.Serve(lis); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()
	
	//создание клиента под grpc
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	
	grpcClient := pb.NewTaskServiceClient(conn)
	
	//запуск агентов в отдельных выражений до переменной COMPUTING_POWER 
	for i := 0; i < countExec; i++ {
		go services.ExecTask(chanTasks, chanResults)
	}
	
	//запуск метода преобразователя в отдельной горутине
	go services.ToSimpleTask(chanPostfixExpression, operationsTime, taskDAO, expressionDAO)
	
	//запуск метода отправления результатов задач серверу в отедльной горутине
	go services.WaiterResults(chanResults, grpcClient)

	//запуск метода приема задач с сервера
	go services.TakerTask(chanTasks, grpcClient)

	//

	r := chi.NewRouter()
	//endpoint для приема выражения с клиента
	
	r.Route("/api/v1", func(r chi.Router) {

		r.Use(auth.CheckUser)
	
		r.Post("/calculate", func(w http.ResponseWriter, r *http.Request) {
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
			
			w.Header().Set("Content-Type", "application/json")
			if errorParse != nil {
				response := map[string]string {
					"error": "not valide expression",
				}
				json.NewEncoder(w).Encode(response)
				log.Println("error: /api/v1/calculate, POST, " + request.Expression)
				w.WriteHeader(422)
				return
			}
			
			userId := r.Context().Value(string("user_id")).(float64)
			
			expression := &models.Expression {
				UserId: int(userId),
				Status: "doing",
				Result: "none",
			}
			
			id, err := expressionDAO.SaveExpression(expression)
			
			postfix := &models.PostfixExpression {
				Id: id,
				Expression: postfixExppression,
			}
			
			chanPostfixExpression <- postfix
			
			log.Println(expression)
			if err != nil {
				log.Println(err)
				return
			}
			
			response := map[string]int {
				"id":id,
			}

			json.NewEncoder(w).Encode(response)
		})	
		
		////endpoint для отправки всех выражений клиенту
		r.Get("/expressions", func(w http.ResponseWriter, r *http.Request) {
			
			userId := r.Context().Value(string("user_id")).(float64)
			
			expressions, _ := expressionDAO.GetAllExpressionByUserId(int(userId))
			
			w.Header().Set("Content-Type", "applciation/json")
			json.NewEncoder(w).Encode(expressions)
			w.WriteHeader(200)
		})
		
		//endpoint для приема выражжения с определенным ID
		r.Get("/expressions/{id}", func(w http.ResponseWriter, r *http.Request) {
			idExpression, _ := strconv.Atoi(chi.URLParam(r, "id"))
			
			userId := r.Context().Value(string("user_id")).(float64)
			
			resultExpression, error := expressionDAO.GetExpressionByIdAndUserId(idExpression, int(userId))
			if error != nil {
				log.Println("error: 404, method: GET")
				return
			}
			
			w.Header().Set("Content-Type", "applciation/json")
			json.NewEncoder(w).Encode(resultExpression)
			w.WriteHeader(200)
		})
	})
	
	
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, req *http.Request) {
				user := models.User{}
				
				err := json.NewDecoder(req.Body).Decode(&user)
				if err != nil {
					log.Println(err)
					return
				}
				
				_, err = userDAO.Save(&user)
				if err != nil {
					log.Println(err)
					return
				}
				
				log.Println("SAVE OK!")
			})
		
		r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
			user := models.User{}
			
			err := json.NewDecoder(req.Body).Decode(&user)
			if err != nil {
				log.Println(err)
				return
			}
			
			userEx, err := userDAO.GetUserByUsername(user.Username)
			if err != nil {
				log.Println(err)
				return
			}
			
			if userEx.Password != user.Password {
				log.Println("error: not eq password")
				return
			}
			
			//can refresh 
			
			token, err := auth.GetTokenString(userEx.Id)
			
			http.SetCookie(w, &http.Cookie {
				Name: "token",
				Value: token,
			})
			
			log.Println(token)
		})
	})
	
	http.ListenAndServe(":8080", r)
}

