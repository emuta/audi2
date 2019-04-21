package service

import (
    "context"

    "github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	_ "audi/pkg/logger"
    "audi/pkg/session/rabbitmq/pubsub/publisher"
    
    pb "audi/account/proto"
	"audi/account/server/repository"
    "audi/account/server/pbtype"
)

type accountServiceServer struct {
	repo   *repository.Repository
	broker *publisher.AppPublisher
}

func NewAccountServiceServer(repo *repository.Repository, broker *publisher.AppPublisher) *accountServiceServer {
	return &accountServiceServer{
		repo:   repo,
		broker: broker,
	}
}

func (s *accountServiceServer) GetRole(ctx context.Context, req *pb.GetRoleReq) (*pb.Role, error) {
    resp, err := s.repo.GetRole(ctx, req.Id, req.Name)
    if err != nil {
        return nil, err
    }
    return pbtype.RoleProto(resp), nil
}


func (s *accountServiceServer) FindRole(ctx context.Context, req *pb.FindRoleReq) (*pb.FindRoleResp, error) {
    args := map[string]interface{}{}

    if req.Id != 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    roles, err := s.repo.FindRole(ctx, args)
    if err != nil {
        return nil, err
    }

    var resp pb.FindRoleResp
    for _, role := range *roles {
        resp.Roles = append(resp.Roles, pbtype.RoleProto(&role))
    }
    return &resp, nil
}


func (s *accountServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {
    user, err := s.repo.CreateUser(ctx, req.Name, req.Password, req.RoleId)
    if err != nil {
        log.WithError(err).Error("Failed to create user")
        return nil, err
    }

    // publish event: user.created
    go s.broker.Publish("user.created", user)


    return pbtype.UserProto(user), nil
}


func (s *accountServiceServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
    user, err := s.repo.GetUser(ctx, req.Id)

    if err != nil {
        log.Error(err)
        return nil, err
    }
    
    return pbtype.UserProto(user), nil
}


func (s *accountServiceServer) ValidateUserPassword(ctx context.Context, req *pb.ValidateUserPasswordReq) (*pb.ValidateUserPasswordResp, error) {
    user, err := s.repo.ValidatePassword(ctx, req.Name, req.Password)
    if err != nil {
        return nil, err
    }
    return &pb.ValidateUserPasswordResp{Id: user.Id, Valid: true}, nil
}

func getFindUserArgs(req *pb.FindUserReq) map[string]interface{} {
    args := map[string]interface{}{}

    if req.Id != 0 {
        args["id"] = req.Id
    }

    if req.Name != "" {
        args["name"] = req.Name
    }

    if req.JoinFrom != nil {
        if t, err := ptypes.Timestamp(req.JoinFrom); err == nil {
            args["join_from"] = t
        }
    }

    if req.JoinTo != nil {
        if t, err := ptypes.Timestamp(req.JoinTo); err == nil {
            args["join_to"] = t
        }
    }

    if req.RoleId != 0 {
        args["role_id"] = req.RoleId
    }

    if req.Active != "" {
        args["active"] = req.Active
    }

    if req.Limit > 0 {
        args["limit"] = req.Limit
    }

    if req.Offset > 0 {
        args["offset"] = req.Offset
    }

    return args
}

func (s *accountServiceServer) FindUser(ctx context.Context, req *pb.FindUserReq) (*pb.FindUserResp, error) {
    users, err := s.repo.FindUser(ctx, getFindUserArgs(req))
    if err != nil {
        return nil, err
    }

    var resp pb.FindUserResp
    for _, user := range *users {
        resp.Users = append(resp.Users, pbtype.UserProto(&user))
    }
    return &resp, nil
}


func (s *accountServiceServer) CountFindUser(ctx context.Context, req *pb.FindUserReq) (*pb.CountFindUserResp, error) {
    total, err := s.repo.CountFindUser(ctx, getFindUserArgs(req))
    if err != nil {
        return nil, err
    }

    return &pb.CountFindUserResp{Total: *total}, nil
}


func (s *accountServiceServer) ChangeUserPassword(ctx context.Context, req *pb.ChangeUserPasswordReq) (*pb.ChangeUserPasswordResp, error) {
    ok, err := s.repo.ChangeUserPassword(ctx, req.Id, req.OldPassword, req.NewPassword)
    if err != nil {
        return nil, err
    }

    // publish event: user.password.changed
    payload := struct {
        OldPassword string
        NewPassword string
    }{
        OldPassword: req.OldPassword,
        NewPassword: req.NewPassword,
    }
    go s.broker.Publish("user.password.changed", payload)

    return &pb.ChangeUserPasswordResp{Value: ok}, nil
}


func (s *accountServiceServer) UpdateUserRole(ctx context.Context, req *pb.UpdateUserRoleReq) (*pb.UpdateUserRoleResp, error) {
    ok, err := s.repo.UpdateUserRole(ctx, req.Id, req.RoleId)
    if err != nil {
        return nil, err
    }

    // publish event: user.role.changed
    payload := struct {
        UserId int64
        RoleId int32
    }{
        UserId: req.Id,
        RoleId: req.RoleId,
    }
    go s.broker.Publish("user.role.changed", payload)

    return &pb.UpdateUserRoleResp{Value: ok}, nil
}


func (s *accountServiceServer) BanUser(ctx context.Context, req *pb.BanUserReq) (*pb.BanUserResp, error) {
    ok, err := s.repo.UpdateUserActive(ctx, req.Id, false)
    if err != nil {
        return nil, err
    }

    // publish event: user.banned
    payload := struct {
        UserId int64
    }{
        UserId: req.Id,
    }
    go s.broker.Publish("user.banned", payload)

    return &pb.BanUserResp{Value: ok}, nil
}


func (s *accountServiceServer) UnBanUser(ctx context.Context, req *pb.UnBanUserReq) (*pb.UnBanUserResp, error) {
    ok, err := s.repo.UpdateUserActive(ctx, req.Id, true)
    if err != nil {
        return nil, err
    }

    // publish event: user.unbanned
    payload := struct {
        UserId int64
    }{
        UserId: req.Id,
    }
    go s.broker.Publish("user.unbanned", payload)

    return &pb.UnBanUserResp{Value: ok}, nil
}


func (s *accountServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
    token, err := s.repo.AuthLogin(ctx, req.Username, req.Password, req.Device, req.Ip, req.Fp)
    if err != nil {
        return nil, err
    }

    // publish event: user.login
    payload := struct {
        Username string
        Password string
        Device   string
        Ip       string
        Fp       string
    }{
        Username: req.Username,
        Password: req.Password,
        Device:   req.Device,
        Ip:       req.Ip,
        Fp:       req.Fp,
    }

    go s.broker.Publish("user.login", payload)

    return &pb.LoginResp{Signature: token.Signature}, nil
}


func (s *accountServiceServer) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutResp, error) {
    if err := s.repo.AuthLogout(ctx, req.Signature); err != nil {
        return nil, err
    }
    // publish event: user.logout
    payload := struct {
        Signature string
    }{
        Signature: req.Signature,
    }
    go s.broker.Publish("user.logout", payload)
    return &pb.LogoutResp{Success: true}, nil
}


func (s *accountServiceServer) GetAuthHistory(ctx context.Context, req *pb.GetAuthHistoryReq) (*pb.AuthHistory, error) {
    resp, err := s.repo.GetAuthHistory(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return pbtype.AuthHistoryProto(resp), nil
}


func getFindAuthHistoryArgs(req *pb.FindAuthHistoryReq) map[string]interface{} {
    args := map[string]interface{}{}

    if req.UserId != 0 {
        args["user_id"] = req.UserId
    }

    if req.Device != "" {
        args["device"] = req.Device
    }

    if req.Ip != "" {
        args["ip"] = req.Ip
    }

    if req.Fp != "" {
        args["fp"] = req.Fp
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

    if req.Limit != 0 {
        args["limit"] = req.Limit
    }

    if req.Offset != 0 {
        args["offset"] = req.Offset
    }

    return args
}


func (s *accountServiceServer) FindAuthHistory(ctx context.Context, req *pb.FindAuthHistoryReq) (*pb.FindAuthHistoryResp, error) {
    results, err := s.repo.FindAuthHistory(ctx, getFindAuthHistoryArgs(req))
    if err != nil {
        return nil, err
    }

    var resp pb.FindAuthHistoryResp
    for _, result := range *results {
        resp.Histories = append(resp.Histories, pbtype.AuthHistoryProto(&result))
    }

    return &resp, nil
}


func (s *accountServiceServer) CountFindAuthHistory(ctx context.Context, req *pb.FindAuthHistoryReq) (*pb.CountFindAuthHistoryResp, error) {
    total, err := s.repo.CountFindAuthHistory(ctx, getFindAuthHistoryArgs(req))
    if err != nil {
        return nil, err
    }
    return &pb.CountFindAuthHistoryResp{Total: *total}, nil
}


func (s *accountServiceServer) GetAuthToken(ctx context.Context, req *pb.GetAuthTokenReq) (*pb.AuthToken, error) {
    resp, err := s.repo.GetAuthToken(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return pbtype.AuthTokenProto(resp), nil
}

func getFindAuthTokenArgs(req *pb.FindAuthTokenReq) map[string]interface{} {
    args := map[string]interface{}{}

    if req.AuthId != 0 {
        args["auth_id"] = req.AuthId
    }

    if req.UserId != 0 {
        args["user_id"] = req.UserId
    }

    if req.UserName != "" {
        args["user_name"] = req.UserName
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

    if req.ExpiredFrom != nil {
        if t, err := ptypes.Timestamp(req.ExpiredFrom); err == nil {
            args["expired_from"] = t
        }
    }

    if req.ExpiredTo != nil {
        if t, err := ptypes.Timestamp(req.ExpiredTo); err == nil {
            args["expired_to"] = t
        }
    }

    if req.Limit != 0 {
        args["limit"] = req.Limit
    }

    if req.Offset != 0 {
        args["offset"] = req.Offset
    }
    return args
}

func (s *accountServiceServer) FindAuthToken(ctx context.Context, req *pb.FindAuthTokenReq) (*pb.FindAuthTokenResp, error) {
    results, err := s.repo.FindAuthToken(ctx, getFindAuthTokenArgs(req))
    if err != nil {
        return nil, err
    }

    var resp pb.FindAuthTokenResp
    for _, result := range *results {
        resp.Tokens = append(resp.Tokens, pbtype.AuthTokenProto(&result))
    }

    return &resp, nil
}


func (s *accountServiceServer) CountFindAuthToken(ctx context.Context, req *pb.FindAuthTokenReq) (*pb.CountFindAuthTokenResp, error) {
    total, err := s.repo.CountFindAuthToken(ctx, getFindAuthTokenArgs(req))
    if err != nil {
        return nil, err
    }
    return &pb.CountFindAuthTokenResp{Total: *total}, nil
}


func (s *accountServiceServer) InactiveAuthToken(ctx context.Context, req *pb.InactiveAuthTokenReq) (*pb.InactiveAuthTokenResp, error) {
    if err := s.repo.InactiveAuthToken(ctx, req.Id); err != nil {
        return nil, err
    }

    // publish event: auth.token.inactived
    payload := struct {
        TokenId int64
    }{
        TokenId: req.Id,
    }
    go s.broker.Publish("auth.token.inactived", payload)
    return &pb.InactiveAuthTokenResp{Success: true}, nil
}


func (s *accountServiceServer) ParseAuthToken(ctx context.Context, req *pb.ParseAuthTokenReq) (*pb.AuthToken, error) {
    resp, err := s.repo.ParseToken(ctx, req.Signature)
    if err != nil {
        return nil, err
    }

    return pbtype.AuthTokenProto(resp), nil
}











































