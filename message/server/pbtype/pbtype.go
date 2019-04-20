package pbtype

import (
	"github.com/golang/protobuf/ptypes"

	pb "audi/message/proto"
	"audi/message/server/model"
)

func MessageProto(msg *model.Message) *pb.Message {
	p := pb.Message{
		Id:     msg.Id,
		Srv:    msg.Srv,
		Event:  msg.Event,
		Active: msg.Active,
	}

	if data, err := msg.Data.MarshalJSON(); err == nil {
		p.Data = data
	}

	if t, err := ptypes.TimestampProto(msg.CreatedAt); err == nil {
		p.CreatedAt = t
	}
	return &p
}

func PublishResultProto(st *model.PublishResult) *pb.PublishResult {
	p := pb.PublishResult{Id: st.Id, MsgId: st.MsgId}
	
	if t, err := ptypes.TimestampProto(st.Ts); err == nil {
		p.Ts = t
	}
	return &p
}