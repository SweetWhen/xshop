package service

import (
	"context"
	mgrpb "realworld/api/mgr/v1"
	"realworld/app/mgr/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
)

var _ mgrpb.UserHTTPServer = &MgrUserSvc{}

// NewMgrUserSvc new a greeter service.
func NewMgrUserSvc(uc *biz.MgrUserUsecase) *MgrUserSvc {
	return &MgrUserSvc{uc: uc}
}

func uInfo2BizUser(info *mgrpb.UserBaseInfo) *biz.User {
	return &biz.User{
		Uid:      info.Uid,
		Account:  info.Account,
		PassWD:   info.Passwd,
		PhoneNum: info.PhoneNum,
		Name:     info.Name,
		Status:   int(info.Status),
	}
}

func bizUser2UInfo(u *biz.User) (res *mgrpb.UserBaseInfo) {
	copier.CopyWithOption(res, u, copier.Option{DeepCopy: true})
	return &mgrpb.UserBaseInfo{
		Uid:      u.Uid,
		Account:  u.Account,
		Passwd:   u.PassWD,
		PhoneNum: u.PhoneNum,
		Name:     u.Name,
		Status:   mgrpb.UserStatus(u.Status),
	}
}

// GreeterService is a greeter service.
type MgrUserSvc struct {
	mgrpb.UnimplementedUserServer

	uc *biz.MgrUserUsecase
}

func (mu *MgrUserSvc) Heartbeat(ctx context.Context, req *mgrpb.HeartbeatRequest) (*mgrpb.HeartbeatReply, error) {
	log.Infof("mgr Heartbeat get a request\n")
	return &mgrpb.HeartbeatReply{Msg: "pong"}, nil
}

// LoginUser implements v1.UserHTTPServer.
func (mu *MgrUserSvc) LoginUser(ctx context.Context, req *mgrpb.LoginUserRequest) (*mgrpb.LoginUserReply, error) {
	resp, err := mu.uc.UserLogin(ctx, req.Account, req.Passwd)
	if err != nil {
		return nil, err
	}
	return &mgrpb.LoginUserReply{Uid: resp.Uid, Name: resp.Name}, nil
}

func (mu *MgrUserSvc) CreateUser(ctx context.Context, req *mgrpb.CreateUserRequest) (resp *mgrpb.CreateUserReply, err error) {
	log.Context(ctx).Debugf("createUser req:%s", req)
	if req.Info == nil {
		return nil, mgrpb.ErrorErrInvalidParam("info is nil")
	}
	u, e := mu.uc.CreateUser(ctx, uInfo2BizUser(req.Info))
	if e != nil {
		return nil, e
	}
	u.PassWD = ""
	return &mgrpb.CreateUserReply{Info: bizUser2UInfo(u)}, nil
}

func (mu *MgrUserSvc) DeleteUser(ctx context.Context, req *mgrpb.DeleteUserRequest) (resp *mgrpb.DeleteUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserSvc) GetUser(ctx context.Context, req *mgrpb.GetUserRequest) (resp *mgrpb.GetUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserSvc) ListUser(ctx context.Context, req *mgrpb.ListUserRequest) (resp *mgrpb.ListUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserSvc) UpdateUser(ctx context.Context, req *mgrpb.UpdateUserRequest) (resp *mgrpb.UpdateUserReply, err error) {
	panic("not implemented") // TODO: Implement
}
