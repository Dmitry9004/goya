package rpc
import (
	pb "goya/project/proto"
	"goya/project/internal/dao"
	"goya/project/internal/models"
	
	"errors"
	"log"
	"time"
	"context"
)

type Server struct {
	service pb.TaskServiceServer
	taskDAO *dao.TaskDAO
}

func NewServer(taskDAO *dao.TaskDAO) *Server {
	return &Server{
		taskDAO: taskDAO,
	}
}

func (s *Server)GetRawTask(ctx context.Context, empty *pb.Empty) (*pb.Task, error) {
		tasks, err := s.taskDAO.GetRawTasks("raw")
		
		if err != nil {
		   log.Println(err)
		   return &pb.Task{}, err
	    }
		if len(tasks) == 0 {
			log.Println("info: empty raw tasks")
			return &pb.Task{}, errors.New("empty raw tasks")
		}
		
		task := tasks[0]
		resTask := &models.Task {
			Id: task.Id,
			Arg1: task.Arg1,
			Arg2: task.Arg2,
			Operation: task.Operation,
			OperationTime: task.OperationTime,
			Status: "process",
		}
		
		rpcTask := &pb.Task {
			Id: int64(task.Id),
			Arg1: task.Arg1,
			Arg2: task.Arg2,
			Operation: task.Operation,
			OperationTime: string(task.OperationTime),
			Status: "process",
		}
		
		s.taskDAO.UpdateTask(resTask)
		
	return rpcTask, nil
}

func (s *Server)PostResultTask(ctx context.Context, task *pb.ResultTask) (*pb.Empty, error) {
		_, err := s.taskDAO.GetTask(int(task.Id))
		
	    if err != nil {
			log.Println(err)
			return &pb.Empty{}, err
	    }
		
		resTask := models.Task{
			Id: int(task.Id),
			Arg1: 0,
			Arg2: 0,
			Result: task.Result,
			Operation: "",
			OperationTime: time.Duration(time.Second),
			Status: "done",
		}
		
		s.taskDAO.UpdateTask(&resTask)
		
	return &pb.Empty{}, nil
}