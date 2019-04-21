package pbtype

import (
	"github.com/golang/protobuf/ptypes"

	pb "audi/account/proto"
	"audi/account/server/model"
)

func RoleProto(m *model.Role) *pb.Role{
	return &pb.Role{Id: m.Id, Name: m.Name}
}

func UserProto(m *model.User) *pb.User {
    p := pb.User{
        Id:     m.Id,
        Name:   m.Name,
        RoleId: m.RoleId,
        Active: m.Active,
    }

    if t, err := ptypes.TimestampProto(m.JoinAt); err == nil {
        p.JoinAt = t
    }

    return &p
}

func AuthHistoryProto(m *model.AuthHistory) *pb.AuthHistory {
	p := pb.AuthHistory{
        Id:     m.Id,
        UserId: m.UserId,
        Device: m.Device,
        Ip:     m.Ip,
        Fp:     m.Fp,
    }

    if t, err := ptypes.TimestampProto(m.CreatedAt); err == nil {
        p.CreatedAt = t
    }

    return &p
}

func AuthTokenProto(m *model.AuthToken) *pb.AuthToken {
    p := pb.AuthToken{
        Id:        m.Id,
        AuthId:    m.AuthId,
        UserId:    m.UserId,
        UserName:  m.UserName,
        Signature: m.Signature,
        Active:    m.Active,
    }

    if t, err := ptypes.TimestampProto(m.CreatedAt); err == nil {
        p.CreatedAt = t
    }

    if t, err := ptypes.TimestampProto(m.ExpiredAt); err == nil {
        p.ExpiredAt = t
    }

    return &p
}