package services

import "goya/project/internal/models"
import "goya/project/internal/dao"

//import "github.com/recoilme/pudge"
import "crypto/rand"
import "math/big"
import "time"
import "strconv"
import "strings"
import "log"

//метод разбивает постифксное выражжение на мелкие "сырые" задачи
//и складывает в базе данных
//после того как все задачи были выполнены, метод сохраняет итоговое представление с результатом
func ToSimpleTask(chanPostfixExpression chan *models.PostfixExpression, operationsTime map[string]time.Duration, taskDAO *dao.TaskDAO, expressionDAO *dao.ExpressionDAO) {

	for {
        select {
            case postfixExppression := <-chanPostfixExpression:
            	resError := ""
				expression := postfixExppression.Expression
                tasksId := models.NewStack[int]()
                st := models.NewStack[string]()
				
				log.Println(expression)
				
                for i := 0; i < len(expression); i++ {
                    if len(expression) > i+2 && isNum(expression[i]) && isNum(expression[i+1]) && !isNum(expression[i+2]) {
        
                        id := GenerateId()
						arg1, _ := strconv.ParseFloat(expression[i], 64)
						arg2, _ := strconv.ParseFloat(expression[i+1], 64)
						
                        task := &models.Task{
                            Id: id,
                            Arg1: arg1,
                            Arg2: arg2,
							Result: 0,
                            Operation: expression[i+2],
                            OperationTime: operationsTime[expression[i]],
							Status: "raw",
                        }
                        
                        tasksId.Push(id)
						
						log.Println(task)
						
						err := taskDAO.SaveTask(task)
						if err != nil {
							log.Println(err)
							resError = "err"
						}
						
                        i += 2
                    } else {
                        if isNum(expression[i]) {
                            st.Push(expression[i] + " " + strconv.Itoa(i))
                        } else {
                            id := GenerateId()
							var arg1 models.ResultTask
							var arg2 models.ResultTask
							
							if st.Len() != 0 { 
								valPre := st.Pop()
								indexString := strings.Split(valPre, " ")[1]
								index, _ := strconv.Atoi(indexString)
								result, _ := strconv.ParseFloat(string(valPre[:len(valPre)-len(indexString)-1]), 64)
								
								if index == i - 1 {
									arg2 = models.ResultTask { Result: result, }
				
									resTask := waitDoneTask(tasksId.Pop(), taskDAO)
									doneTask := models.ResultTask{
										Id: resTask.Id,
										Result: resTask.Result,
									}
									arg1 = doneTask
								} else {
									arg1 = models.ResultTask { Result: result, }
									
									resTask := waitDoneTask(tasksId.Pop(), taskDAO)
							
									doneTask := models.ResultTask{
										Id: resTask.Id,
										Result: resTask.Result,
									}
									arg2 = doneTask
									
								}
                            } else if tasksId.Len() != 0 {
								resTask := waitDoneTask(tasksId.Pop(), taskDAO)
									doneTask := models.ResultTask{
										Id: resTask.Id,
										Result: resTask.Result,
									}
								arg2 = doneTask
							
								resTask = waitDoneTask(tasksId.Pop(), taskDAO)
									doneTask = models.ResultTask{
										Id: resTask.Id,
										Result: resTask.Result,
									}
								arg1 = doneTask
									
							} else {
								resError = "err"
								break
							}
                            task := &models.Task {
                                Id: id,
                                Arg1: arg1.Result,
                                Arg2: arg2.Result,
								Result: 0,
                                Operation: expression[i],
                                OperationTime: operationsTime[expression[i]],
								Status: "raw",
                            }
							if task.Operation == "/" && arg2.Result == 0 {
								resError = "err"
								break
							}

                            err := taskDAO.SaveTask(task)
							if err != nil {
								log.Println(err)
								resError = "err"
								break
							}
							tasksId.Push(id)
						}                        
                    }
                }

                if len(resError) != 0 {
					//need save expression with status 'not valid data'
					
					origExpression, _ := expressionDAO.GetExpressionById(postfixExppression.Id)
					
					failExpression := &models.Expression {
						Id: origExpression.Id,
						UserId: origExpression.UserId,
						Result: "",
						Status: "not valid expression",
					}
					expressionDAO.UpdateExpression(failExpression)
                	continue
                }
				
				var taskResult models.Task
				id := tasksId.Peek()
				if tasksId.Len() != 0 {
					for {
						task, err := taskDAO.GetTask(id)
						if err == nil && task.Status == "done" {
							taskResult = task
							break
						}
					}
					tasksId.Pop()
				} 
		
				result := strconv.FormatFloat(taskResult.Result, 'f', -1, 64)
				
				resultExpression := &models.Expression {
					Id:postfixExppression.Id,
					Status: "Done",
					Result: result,
				}
				
				log.Println(resultExpression)
				
				errSave := expressionDAO.UpdateExpression(resultExpression)
				
				if errSave != nil {
					log.Println(errSave)
				}
		}
    }
}

func waitDoneTask(taskId int, taskDAO *dao.TaskDAO) models.Task {
    for {
		task, err := taskDAO.GetTask(taskId)
		if err == nil && task.Status == "done" {
			return task
		}
    }
	return models.Task{}
}

//метод генерации ID для представлений с использованием 
//crypto/rand and math/big 
func GenerateId() int {
    val, _ := rand.Int(rand.Reader, big.NewInt(10000000))
	return int(val.Uint64())
}
