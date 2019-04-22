package service

import (
	"context"
    "time"

    "github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	_ "audi/pkg/logger"
	"audi/pkg/session/rabbitmq/pubsub/publisher"
	pb "audi/lottery/ssc/cqssc/proto"
	"audi/lottery/ssc/cqssc/server/repository"
	"audi/lottery/ssc/cqssc/server/pbtype"
)

type cqsscServiceServer struct {
		repo  *repository.Repository
		broker *publisher.AppPublisher
}

func NewCqsscServiceServer(repo *repository.Repository, broker *publisher.AppPublisher) * cqsscServiceServer {
	log.Info("Service loaded")
	return &cqsscServiceServer{
		repo:  repo,
		broker: broker,
	}
}

func (s *cqsscServiceServer) GetConfig(ctx context.Context, req *pb.GetConfigReq) (*pb.Config, error) {
    resp, err := s.repo.GetConfig(ctx)
    if err != nil {
        log.WithError(err).Error("Failed to get config")
        return nil, err
    }
    return pbtype.ConfigProto(resp), nil
}

func (s *cqsscServiceServer) UpdateConfig(ctx context.Context, req *pb.UpdateConfigReq) (*pb.Config, error) {
    resp, err := s.repo.UpdateConfig(ctx, req.Odds, req.Comm)
    if err != nil {
        log.WithError(err).Error("Failed to update config")
        return nil, err
    }
    go s.broker.Publish("config.updated", resp)
    return pbtype.ConfigProto(resp), nil
}

func (s *cqsscServiceServer) GetUnit(ctx context.Context, req *pb.GetUnitReq) (*pb.Unit, error) {
    resp, err := s.repo.GetUnit(ctx, req.Id)
    if err != nil {
        log.WithError(err).Error("Failed to get unit")
        return nil, err
    }
    return pbtype.UnitProto(resp), nil
}

func (s *cqsscServiceServer) FindUnit(ctx context.Context, req *pb.FindUnitReq) (*pb.FindUnitResp, error) {
    args := map[string]interface{}{}
    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    if req.Value > 0 {
        args["value"] = req.Value
    }

    units, err := s.repo.FindUnit(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindUnitResp
    for _, unit := range *units {
        resp.Units = append(resp.Units, pbtype.UnitProto(&unit))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) GetCatg(ctx context.Context, req *pb.GetCatgReq) (*pb.Catg, error) {
    resp, err := s.repo.GetCatg(ctx, req.Id)
    if err != nil {
        log.WithError(err).Error("Failed to get unit")
        return nil, err
    }
    return pbtype.CatgProto(resp), nil
}

func (s *cqsscServiceServer) FindCatg(ctx context.Context, req *pb.FindCatgReq) (*pb.FindCatgResp, error) {
    args := map[string]interface{}{}
    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    if req.Tag != "" {
        args["tag"] = req.Tag
    }

    if req.Pref != "" {
        args["pref"] = req.Pref
    }

    cats, err := s.repo.FindCatg(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindCatgResp
    for _, catg := range *cats {
        resp.Catgs = append(resp.Catgs, pbtype.CatgProto(&catg))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.Group, error) {
    resp, err := s.repo.GetGroup(ctx, req.Id)
    if err != nil {
        log.WithError(err).Error("Failed to get unit")
        return nil, err
    }
    return pbtype.GroupProto(resp), nil
}

func (s *cqsscServiceServer) FindGroup(ctx context.Context, req *pb.FindGroupReq) (*pb.FindGroupResp, error) {
    args := map[string]interface{}{}
    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    if req.Tag != "" {
        args["tag"] = req.Tag
    }

    groups, err := s.repo.FindGroup(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindGroupResp
    for _, group := range *groups {
        resp.Groups = append(resp.Groups, pbtype.GroupProto(&group))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) GetPlay(ctx context.Context, req *pb.GetPlayReq) (*pb.Play, error) {
    resp, err := s.repo.GetPlay(ctx, req.Id)
    if err != nil {
        log.WithError(err).Error("Failed to get unit")
        return nil, err
    }
    return pbtype.PlayProto(resp), nil
}

func (s *cqsscServiceServer) FindPlay(ctx context.Context, req *pb.FindPlayReq) (*pb.FindPlayResp, error) {
    args := map[string]interface{}{}
    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    if req.Tag != "" {
        args["tag"] = req.Tag
    }

    if req.Pref != "" {
        args["pref"] = req.Pref
    }

    if req.Active != "" {
        args["active"] = req.Active
    }

    if req.Pr > 0 {
        args["pr"] = req.Pr
    }

    if req.CatgId > 0 {
        args["catg_id"] = req.CatgId
    }

    if req.GroupId > 0 {
        args["group_id"] = req.GroupId
    }

    if req.UnitId > 0 {
        args["unit_id"] = req.UnitId
    }

    items, err := s.repo.FindPlay(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindPlayResp
    for _, item := range *items {
        resp.Plays = append(resp.Plays, pbtype.PlayProto(&item))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) UpdatePlay(ctx context.Context, req *pb.UpdatePlayReq) (*pb.Play, error) {
    resp, err := s.repo.UpdatePlay(ctx, req.Id, req.Pref, req.Active, req.Units)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("play.updated", resp)
    return pbtype.PlayProto(resp), nil
}

func (s *cqsscServiceServer) CreateTerm(ctx context.Context, req *pb.CreateTermReq) (*pb.Term, error) {
    var st, et time.Time
    var err error

    if st, err = ptypes.Timestamp(req.StartFrom); err != nil {
        return nil, err
    }

    if et, err = ptypes.Timestamp(req.EndTo); err != nil {
        return nil, err
    }

    resp, err := s.repo.CreateTerm(ctx, req.Id, st, et)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("term.created", resp)
    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) GetTerm(ctx context.Context, req *pb.GetTermReq) (*pb.Term, error) {
    resp, err := s.repo.GetTerm(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) FindTerm(ctx context.Context, req *pb.FindTermReq) (*pb.FindTermResp, error) {
    args := map[string]interface{}{}

    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.StartFrom != nil {
        if t, err := ptypes.Timestamp(req.StartFrom); err != nil {
            args["start_from"] = t
        }
    }

    if req.EndTo != nil {
        if t, err := ptypes.Timestamp(req.EndTo); err != nil {
            args["end_to"] = t
        }
    }

    if req.Limit > 0 {
        args["limit"] = req.Limit
    }

    if req.Offset > 0 {
        args["offset"] = req.Offset
    }

    if req.OrderBy != "" {
        args["order_by"] = req.OrderBy
    }

    terms, err := s.repo.FindTerm(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindTermResp
    for _, period := range *terms {
        resp.Terms = append(resp.Terms, pbtype.TermProto(&period))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) CountFindTerm(ctx context.Context, req *pb.FindTermReq) (*pb.CountFindTermResp, error) {
    args := map[string]interface{}{}

    if req.Id > 0 {
        args["id"] = req.Id
    }

    if req.StartFrom != nil {
        if t, err := ptypes.Timestamp(req.StartFrom); err != nil {
            args["start_from"] = t
        }
    }

    if req.EndTo != nil {
        if t, err := ptypes.Timestamp(req.EndTo); err != nil {
            args["end_to"] = t
        }
    }

    if req.Limit > 0 {
        args["limit"] = req.Limit
    }

    if req.Offset > 0 {
        args["offset"] = req.Offset
    }

    if req.OrderBy != "" {
        args["order_by"] = req.OrderBy
    }

    total, err := s.repo.CountFindTerm(ctx, args)
    if err != nil {
        return nil, err
    }
    return &pb.CountFindTermResp{Total: *total}, nil
}

func (s *cqsscServiceServer) SetTermResult(ctx context.Context, req *pb.SeTermtResultReq) (*pb.Term, error) {
    resp, err := s.repo.SetTermResult(ctx, req.Id, req.Codes)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("term.result.created", resp)
    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) UpdateTermResult(ctx context.Context, req *pb.UpdateTermResultReq) (*pb.Term, error) {
    resp, err := s.repo.SetTermResult(ctx, req.Id, req.Codes)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("term.result.updated", resp)
    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) SettleTerm(ctx context.Context, req *pb.SettleTermReq) (*pb.Term, error) {
    resp, err := s.repo.SettleTerm(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("term.settled", resp)
    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) RevokeTerm(ctx context.Context, req *pb.RevokeTermReq) (*pb.Term, error) {
    resp, err := s.repo.RevokeTerm(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    go s.broker.Publish("term.revoked", resp)
    return pbtype.TermProto(resp), nil
}

func (s *cqsscServiceServer) CreateBet(ctx context.Context, req *pb.CreateBetReq) (*pb.Bet, error) {
    var plans []map[string]interface{}
    for _, plan := range req.Plans {
        plans = append(plans, map[string]interface{}{
            "term_id": plan.TermId,
            "times": plan.Times,
            "qty": plan.Qty,
            "payment": plan.Payment,
            "bonus": plan.Bonus,
            "rebate": plan.Rebate,
        })
    }
    log.Info(plans)
    resp, err := s.repo.CreateBet(ctx, 
        req.UserId, req.Odds, req.PlayId, req.UnitId, req.Comm, req.ChaseStop, req.Codes, plans)
    if err != nil {
        return nil, err
    }
    // should publish event
    go s.broker.Publish("bet.created", resp)
    return pbtype.BetProto(resp), nil
}

func (s *cqsscServiceServer) GetBet(ctx context.Context, req *pb.GetBetReq) (*pb.Bet, error) {
    resp, err := s.repo.GetBet(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return pbtype.BetProto(resp), nil
}

func getFindBetRequestArgument(req *pb.FindBetReq) map[string]interface{} {
    args := map[string]interface{}{}

    if req.UserId > 0 {
        args["user_id"] = req.UserId
    }

    if req.PlayId > 0 {
        args["play_id"] = req.PlayId
    }

    if req.Code != "" {
        args["code"] = req.Code
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

    if req.TermId > 0 {
        args["term_id"] = req.TermId
    }

    if req.Win != "" {
        args["win"] = req.Win
    }

    if req.Settled != "" {
        args["settled"] = req.Settled
    }

    if req.Revoked != "" {
        args["revoked"] = req.Revoked
    }

    if req.Limit > 0 {
        args["limit"] = req.Limit
    }

    if req.Offset > 0 {
        args["offset"] = req.Offset
    }

    return args
}

func (s *cqsscServiceServer) FindBet(ctx context.Context, req *pb.FindBetReq) (*pb.FindBetResp, error) {
    args := getFindBetRequestArgument(req)

    bs, err := s.repo.FindBet(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindBetResp
    for _, b := range *bs {
        resp.Bets = append(resp.Bets, pbtype.BetProto(&b))
    }

    return &resp, nil
}

func (s *cqsscServiceServer) CountFindBet(ctx context.Context, req *pb.FindBetReq) (*pb.CountFindBetResp, error) {
    args := getFindBetRequestArgument(req)

    total, err := s.repo.CountFindBet(ctx, args)
    if err != nil {
        return nil, err
    }
   
    return &pb.CountFindBetResp{Total: *total}, nil
}

func (s *cqsscServiceServer) GetBetStats(ctx context.Context, req *pb.GetBetStatsReq) (*pb.GetBetStatsResp, error) {
    stats, err := s.repo.GetBetStats(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    var resp pb.GetBetStatsResp
    for _, stat := range *stats {
        resp.Stats = append(resp.Stats, pbtype.BetPlanStatsProto(&stat))
    }
    return &resp, nil
}

func (s *cqsscServiceServer) GetBetPlan(ctx context.Context, req *pb.GetBetPlanReq) (*pb.BetPlan, error) {
    resp, err := s.repo.GetBetPlan(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return pbtype.BetPlanProto(resp), nil
}

func (s *cqsscServiceServer) GetBetPlanStats(ctx context.Context, req *pb.GetBetPlanStatsReq) (*pb.BetPlanStats, error) {
    resp, err := s.repo.GetBetPlanStats(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return pbtype.BetPlanStatsProto(resp), nil
}

func (s *cqsscServiceServer) SettleBetPlan(ctx context.Context, req *pb.SettleBetPlanReq) (*pb.SettleBetPlanResp, error) {
    stats, err := s.repo.SettleBetPlan(ctx, req.Id, req.Win)
    if err != nil {
        return nil, err
    }

    log.WithFields(log.Fields{
        "stats": stats,
    }).Info("Bet plan settled")
    
    // should publish event
    // assign bonus
    // assign rebate

    return &pb.SettleBetPlanResp{Success: true}, nil
}

func (s *cqsscServiceServer) RevokeBetPlan(ctx context.Context, req *pb.RevokeBetPlanReq) (*pb.RevokeBetPlanResp, error) {
    stats, err := s.repo.RevokeBetPlan(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    log.WithFields(log.Fields{
        "stats": stats,
    }).Info("Bet plan revoked")

    // should publish evet
    // refund payment
    // revoke bonus
    // revoke rebate

    return &pb.RevokeBetPlanResp{Success: true}, nil
}