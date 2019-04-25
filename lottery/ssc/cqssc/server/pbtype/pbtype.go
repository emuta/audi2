package pbtype

import (
	"github.com/golang/protobuf/ptypes"

	pb "audi/lottery/ssc/cqssc/proto"
	"audi/lottery/ssc/cqssc/server/model"
)

func ConfigProto(m *model.Config) *pb.Config {
	return &pb.Config{
		Name:   m.Name,
		Tag:    m.Tag,
		Odds:   m.Odds,
		Comm:   m.Comm,
		Price:  m.Price,
		Active: m.Active,
	}
}

func UnitProto(m *model.Unit) *pb.Unit {
	return &pb.Unit{Id: m.Id, Name: m.Name, Value: m.Value}
}

func CatgProto(m *model.Catg) *pb.Catg {
	return &pb.Catg{
		Id:   m.Id,
		Name: m.Name,
		Tag:  m.Tag,
		Pref: m.Pref,
	}
}

func GroupProto(m *model.Group) *pb.Group {
	return &pb.Group{Id: m.Id, Name: m.Name}
}

func PlayProto(m *model.Play) *pb.Play {
	return &pb.Play{
		Id:      m.Id,
		Name:    m.Name,
		Tag:     m.Tag,
		Pref:    m.Pref,
		Active:  m.Active,
		Pr:      m.Pr,
		CatgId:  m.CatgId,
		GroupId: m.GroupId,
		Units:   m.Units,
	}
}


func TermProto(m *model.Term) *pb.Term {
	p := pb.Term{Id: m.Id, Codes: m.Codes}

	if t, err := ptypes.TimestampProto(m.StartFrom); err != nil {
		p.StartFrom = t
	}

	if t, err := ptypes.TimestampProto(m.EndTo); err == nil {
		p.EndTo = t
	}

	if m.OpenedAt.Valid {
		if t, err := ptypes.TimestampProto(m.OpenedAt.Time); err == nil {
			p.OpenedAt = t
		}
	}

	if m.SettledAt.Valid {
		if t, err := ptypes.TimestampProto(m.SettledAt.Time); err == nil {
			p.SettledAt = t
		}
	}

	if m.RevokedAt.Valid {
		if t, err := ptypes.TimestampProto(m.RevokedAt.Time); err == nil {
			p.RevokedAt = t
		}
	}

	return &p
}

func BetProto(m *model.Bet) *pb.Bet {
	p := pb.Bet{
		Id:        m.Id,
		UserId:    m.UserId,
		Odds:      m.Odds,
		PlayId:    m.PlayId,
		UnitId:    m.UnitId,
		Comm:      m.Comm,
		ChaseStop: m.ChaseStop,
		Codes:     m.Codes,
	}

	if t, err := ptypes.TimestampProto(m.CreatedAt); err == nil {
		p.CreatedAt = t
	}

	return &p
}

func BetPlanProto(m *model.BetPlan) *pb.BetPlan {
	p := pb.BetPlan{
		Id:      m.Id,
		BetId:   m.BetId,
		TermId:  m.TermId,
		Times:   m.Times,
		Qty:     m.Qty,
		Payment: m.Payment,
		Rebate:  m.Rebate,
		Bonus:   m.Bonus,
	}

	return &p
}

func BetPlanStatsProto(m *model.BetPlanStats) *pb.BetPlanStats {
	p := pb.BetPlanStats{
		Id:      m.Id,
		BetId:   m.BetId,
		Settled: m.Settled,
		Revoked: m.Revoked,
		Payment: m.Payment,
		Rebate:  m.Rebate,
		Bonus:   m.Bonus,
		Win:     m.Win,
	}

	if m.SettledAt.Valid {
		if t, err := ptypes.TimestampProto(m.SettledAt.Time); err == nil {
			p.SettledAt = t
		}
	}

	if m.RevokedAt.Valid {
		if t, err := ptypes.TimestampProto(m.RevokedAt.Time); err == nil {
			p.RevokedAt = t
		}
	}

	return &p
}