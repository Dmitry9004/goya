package services

import "goya/project/internal/models"
import "github.com/recoilme/pudge"
import "crypto/rand"
import "math/big"
import "time"
import "strconv"
import "strings"
import "log"

//метод разбивает постифксное выражжение на мелкие "сырые" задачи
//и складывает в базе данных
//после того как все задачи были выполнены, метод сохраняет итоговое представление с результатом
func ToSimpleTask(chanPostfixExpression chan *models.PostfixExpression, operationsTime map[string]time.Duration) {
    db, _ := pudge.Open("../expressions", nil)
    dbTasksDone, _ := pudge.Open("../tasks/done", nil)
	dbTasksRaw, _ := pudge.Open("../tasks/raw", nil)
	
	defer db.Close()
	defer dbTasksDone.Close()
	defer dbTasksRaw.Close()

	for {
        select {
            case postfixExppression := <-chanPostfixExpression:
            	err := ""
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
                            Operation: expression[i+2],
                            Operation_time: operationsTime[expression[i]],
                        }
                        
                        tasksId.Push(id)
                        dbTasksRaw.Set(id, task)
												
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
									
									waitDoneTask(dbTasksDone, tasksId.Peek())
									dbTasksDone.Get(tasksId.Pop(), &arg1)
								} else {
									arg1 = models.ResultTask { Result: result, }
									
									waitDoneTask(dbTasksDone, tasksId.Peek())
									dbTasksDone.Get(tasksId.Pop(), &arg2)
								}
                            } else {
                                waitDoneTask(dbTasksDone, tasksId.Peek())
								dbTasksDone.Get(tasksId.Pop(), &arg2)
                                
                                waitDoneTask(dbTasksDone, tasksId.Peek())
								dbTasksDone.Get(tasksId.Pop(), &arg1)
							}
                            
                            task := &models.Task {
                                Id: id,
                                Arg1: arg1.Result,
                                Arg2: arg2.Result,
                                Operation: expression[i],
                                Operation_time: operationsTime[expression[i]],
                            }
							if task.Operation == "/" && arg2.Result == 0 {
								db.Delete(postfixExppression.Id)

								resultExpression := &models.Expression {
									Id:postfixExppression.Id,
									Status: "Not valide expression",
									Result: "...",
								}

								db.Set(postfixExppression.Id, resultExpression)
								err = "err"
								break
							}

                            dbTasksRaw.Set(id, task)
							tasksId.Push(id)
						}                        
                    }
                }

                if len(err) != 0 {
                	continue
                }
				
				for {
					if t, _ := dbTasksDone.Has(tasksId.Peek()); t {
						break
					}
				}
		
				db.Delete(postfixExppression.Id)
				var task models.ResultTask
				dbTasksDone.Get(tasksId.Pop(), &task)
			
				result := strconv.FormatFloat(task.Result, 'f', -1, 64)
				
				resultExpression := &models.Expression {
					Id:postfixExppression.Id,
					Status: "Done",
					Result: result,
				}

				db.Set(postfixExppression.Id, resultExpression)
        }
    }
}

func waitDoneTask(dbTasksDone *pudge.Db, taskId int) {
    for {
        if t, _ := dbTasksDone.Has(taskId); t {
            break;
        }
    }
}

//метод генерации ID для представлений с использованием 
//crypto/rand and math/big 
func GenerateId() int {
    val, _ := rand.Int(rand.Reader, big.NewInt(10000000))
	return int(val.Uint64())
}
