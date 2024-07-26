package services

import (
	"goya/project/internal/models"
	pb "goya/project/proto"
	
	"regexp"
	"errors"
	"log"
	"time" 
	"context"
)

//метод проверки знака на число
func isNum(str string) bool {
    t, _ := regexp.MatchString("[0-9]", str) 
    return t
}

//метод преобразования выражения в постфиксное
func ToPostfix(str string) ([]string, error) {
    //операции
    prt := map[string]int {
        "*": 3,
        "/": 3,
        "+": 2,
        "-": 2,
        "(": 1,
        ")": 1,
    }
    
    //stack операций
    st := models.NewStack[string]()
    res := []string{}
    
    for i := 0; i < len(str); i++ {
        ch := string(str[i])
        
        if ch == "-" && (i-1 >= 0 && string(str[i-1]) == "(" || st.Len() == 0) {
            num, index := getFullNum(str, i+1)
            res = append(res, "-"+num)
            i = index
            continue
        }
        
        if isNum(ch) {
            num, index := getFullNum(str, i)
            res = append(res, num)
            i = index
        } else if prt[ch] != 0 {
            if ch == "(" {
                st.Push("(")
                continue
            }
            
            if ch == ")" {
                for st.Len() != 0 && st.Peek() != "(" {
                    res = append(res, st.Pop())
                }
                if st.Len() != 0 {
                    st.Pop()
                }
                continue
            }
            
            if prt[st.Peek()] >= prt[ch] {
                for prt[st.Peek()] >= prt[ch] {
                  
                    res = append(res, st.Pop())
                }
                
                st.Push(ch)
            } else {
                st.Push(ch)
            }
        } else { 
            return []string{}, errors.New("not find operation")
        }
    }
    
    for st.Len() != 0 {
    	if (st.Peek() == "(" || st.Peek() == ")") {
    		return []string{}, errors.New("not valide expression")
    	}

        res = append(res, st.Pop())
    }
    
    return res, nil
}


//метод возвращает все число и последний индекс числа
func getFullNum(str string, i int) (string, int) {
    currentNum := ""
    lastIndex := len(str)
    for index := i; index < len(str); index++ {
        ch := string(str[index])
        if isNum, _ := regexp.MatchString("[0-9]", ch); isNum {
            currentNum += ch
        } else { 
            lastIndex = index-1
            break
        }
    }
    
    return currentNum, lastIndex
}


//метод выолняет задачи и отпрвляет результат в канал
func ExecTask(chanTasks chan models.Task, chanResults chan *models.ResultTask) {
	for {
		select {
			case task := <- chanTasks:
				timer := time.NewTimer(task.OperationTime)
				select {
					case <-timer.C:	
						result, err := calculate(task) 
						if err != nil {
							log.Print(err)					
						}
						
						responseResult := &models.ResultTask {
							Id: task.Id,
							Result: result,
						}
						
						chanResults <- responseResult
				}
		}
	}
}

//метод ждет результат в канале и после отправляет на сервер
func WaiterResults(chanResults chan *models.ResultTask, client pb.TaskServiceClient) {
	for {
		select {
			case result := <- chanResults:
				
				messageResultTask := &pb.ResultTask {
					Id: int64(result.Id),
					Result: result.Result,
				}
				
				//вызов grpc метода
				_, err := client.PostResultTask(context.TODO(), messageResultTask)
				
				if err != nil {
					log.Println("error from waiterResults")
					continue
				}
		}
	}
}

//метод запрашивает с сервера задачи и отправляет их в канал
func TakerTask(chanTasks chan models.Task, client pb.TaskServiceClient) {
	ticker := time.NewTicker(time.Second)
	
	for {
		select {
			case <- ticker.C: 
				
				rpcTask, err := client.GetRawTask(context.TODO(), &pb.Empty{})
				
				if err != nil {
					log.Println(err)
					continue
				}
				
				time, _ := time.ParseDuration(rpcTask.OperationTime)
				task := models.Task {
					Id: int(rpcTask.Id),
					ExpressionId: int(rpcTask.ExpressionId),
					Arg1: rpcTask.Arg1,
					Arg2: rpcTask.Arg2,
					Result: rpcTask.Result,
					Operation: rpcTask.Operation,
					OperationTime: time,
					Status: rpcTask.Status,
				}
				
				chanTasks <- task
		}
	}
}

//метод выполнения простейших операций
func calculate(task models.Task) (float64, error) {	
	result := 0.0
	
	switch(task.Operation) { 
		case "+": 
			result = task.Arg1 + task.Arg2
		case "-":
			result = task.Arg1 - task.Arg2
		case "*":
			result = task.Arg1 * task.Arg2
		case "/":
			if task.Arg2 == 0 {
				return 0, errors.New("arg2 = 0")
			}
			result = task.Arg1 / task.Arg2
		default:
			log.Print("not find operation: " + string(task.Operation) + "id: " + string(task.Id))
			return 0, errors.New("not find operation") 
	}

	return result, nil
}
