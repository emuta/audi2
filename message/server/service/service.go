package service

import (
	"context"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	_ "audi/pkg/logger"
	"audi/pkg/session/rabbitmq/pubsub/producer"
    
    pb "audi/message/proto"
	"audi/message/server/repository"
    "audi/message/server/pbtype"
)

type messageServiceServer struct {
	repo   *repository.Repository
	broker *producer.Producer
}

func NewMessageServiceServer(repo *repository.Repository, broker *producer.Producer) *messageServiceServer {
	return &messageServiceServer{
		repo:   repo,
		broker: broker,
	}
}

func (s *messageServiceServer) Create(ctx context.Context, req *pb.CreateReq) (*pb.Message, error) {
    resp, err := s.repo.CreateMessage(ctx, req.Srv, req.Event, req.Data)
	if err != nil {
		return nil, err
	}

	return pbtype.MessageProto(resp), nil
}


func (s *messageServiceServer) Get(ctx context.Context, req *pb.GetReq) (*pb.Message, error) {
    resp, err := s.repo.GetMessage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return pbtype.MessageProto(resp), nil
}


func (s *messageServiceServer) Delete(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteResp, error) {
    if err := s.repo.DeleteMessage(ctx, req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteResp{Value: true}, nil
}

func getFindMessageArgs(req *pb.FindReq) map[string]interface{} {
	args := map[string]interface{}{}

    if req.Id > 0 {
    	args["id"] = req.Id
    }

    if req.Srv != "" {
    	args["srv"] = req.Srv
    }

    if req.Event != "" {
    	args["event"] = req.Event
    }

    if req.Data != "" {
    	args["data"] = req.Data
    }

    if req.CreatedFrom != nil {
    	if t, err := ptypes.Timestamp(req.CreatedFrom); err == nil {
    		args["created_from"] = t
    	}
    }

    if req.CreatedTo != nil {
    	if t, err := ptypes.Timestamp(req.CreatedTo); err == nil {
    		args["created_to"] = t
    	}
    }

    if req.Active != "" {
    	if v, err := strconv.ParseBool(req.Active); err == nil {
    		args["active"] = v
    	}
    }

    if req.Limit > 0 {
    	args["limit"] = req.Limit
    }

    if req.Offset > 0 {
    	args["offset"] = req.Offset
    }

    return args
}

func (s *messageServiceServer) Find(ctx context.Context, req *pb.FindReq) (*pb.FindResp, error) {
	args := getFindMessageArgs(req)
    msgs, err := s.repo.FindMessage(ctx, args)
	if err != nil {
		return nil, err
	}

	var resp pb.FindResp
	for _, msg := range *msgs {
		resp.Messages = append(resp.Messages, pbtype.MessageProto(&msg))
	}
	return &resp, nil
}


func (s *messageServiceServer) CountFind(ctx context.Context, req *pb.FindReq) (*pb.CountFindResp, error) {
    args := getFindMessageArgs(req)
    total, err := s.repo.CountFindMessage(ctx, args)
	if err != nil {
		return nil, err
	}
	return &pb.CountFindResp{Total: *total}, nil
}


func (s *messageServiceServer) Publish(ctx context.Context, req *pb.PublishReq) (*pb.PublishResult, error) {
    msg, err := s.repo.GetMessage(ctx, req.MsgId)
	if err != nil {
		return nil, err
	}

	// publish message
	ch := make(chan error)
	go func(){
		defer close(ch)

		data, err := msg.Data.MarshalJSON()
		if err != nil {
			ch <- err
			return
		}
		payload := producer.NewMessage(msg.Srv, msg.Event, msg.Id, data)
		ch <- s.broker.Publish(payload)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	result, err := s.repo.CreatePublish(ctx, req.MsgId)
	if err != nil {
		return nil, err
	}

	return pbtype.PublishResultProto(result), nil
}


func (s *messageServiceServer) FindPublish(ctx context.Context, req *pb.FindPublishReq) (*pb.FindPublishResp, error) {
    results, err := s.repo.FindPublish(ctx, req.Id, req.MsgId)
    if err != nil {
    	return nil, err
    }

    var resp pb.FindPublishResp
    for _, result := range *results {
    	resp.Results = append(resp.Results, pbtype.PublishResultProto(&result))
    }
    return &resp, nil
}

















