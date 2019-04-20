package service

import (
	"context"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/golang/protobuf/ptypes"
	// log "github.com/sirupsen/logrus"

	pb "audi/protobuf/message"
	"audi/message/pubsub/producer"
	"audi/message/server/pbtype"
	"audi/message/server/repo"
)

var ExchangePubSubName string

type messageServiceServer struct {
	db       *gorm.DB
	producer *producer.Producer
	repo     *repo.Repo
}

func NewMessageServiceServer(db *gorm.DB, p *producer.Producer) *messageServiceServer {
	return &messageServiceServer{
		producer: p, 
		repo:     repo.NewRepo(db),
	}
}

func (s *messageServiceServer) CreateMessage(ctx context.Context, req *pb.CreateMessageReq) (*pb.Message, error) {
	msg, err := s.repo.CreateMessage(ctx, req.Srv, req.Event, postgres.Jsonb{req.Data})
	if err != nil {
		return nil, err
	}

	return pbtype.MessageProto(msg), nil
}

func (s *messageServiceServer) GetMessage(ctx context.Context, req *pb.GetMessageReq) (*pb.Message, error) {
	msg, err := s.repo.GetMessage(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return pbtype.MessageProto(msg), nil
}

func (s *messageServiceServer) DeleteMessage(ctx context.Context, req *pb.DeleteMessageReq) (*pb.DeleteMessageResp, error) {
	if err := s.repo.DeleteMessage(ctx, req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteMessageResp{Value: true}, nil
}

func (s *messageServiceServer) FindMessage(ctx context.Context, req *pb.FindMessageReq) (*pb.FindMessageResp, error) {
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

	msgs, err := s.repo.FindMessage(ctx, args)
	if err != nil {
		return nil, err
	}

	var resp pb.FindMessageResp
	for _, msg := range *msgs {
		resp.Messages = append(resp.Messages, pbtype.MessageProto(&msg))
	}
	return &resp, nil
}

func (s *messageServiceServer) GetPublishHistory(ctx context.Context, req *pb.PublishMessageReq) (*pb.PublishMessageResp, error) {
	ps, err := s.repo.GetPublishes(ctx, req.Id, 10, 0)
	if err != nil {
		return nil, err
	}

	var resp pb.PublishMessageResp
	for _, p := range *ps {
		resp.Publishes = append(resp.Publishes, pbtype.PublishProto(&p))
	}

	return &resp, nil
}

func (s *messageServiceServer) PublishMessage(ctx context.Context, req *pb.PublishMessageReq) (*pb.PublishMessageResp, error) {
	// get message
	msg, err := s.repo.GetMessage(ctx, req.Id)
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

		// routingkye, appId, type, messageId, body
		ch <- s.producer.Publish("", msg.Srv, msg.Event, strconv.FormatInt(msg.Id, 10), data)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <- ch:
		if err != nil {
			return nil, err
		}
	}

	// create publish history
	_, err = s.repo.CreatePublish(ctx, req.Id)
	if err != nil {
		return nil, err
	}


	// return all last publish history
	publishes, err := s.repo.GetPublishes(ctx, req.Id, 10, 0)
	if err != nil {
		return nil, err
	}

	var resp pb.PublishMessageResp
	for _, p := range *publishes {
		resp.Publishes = append(resp.Publishes, pbtype.PublishProto(&p))
	}

	return &resp, nil
}







