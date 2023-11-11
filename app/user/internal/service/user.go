package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	userpb "realworld/api/user/v1"
	"realworld/app/user/internal/biz"
)

type UserService struct {
	userpb.UnimplementedUserServer

	log log.Logger

	ubiz *biz.UserUsecase
}

var _ userpb.UserServer = &UserService{}

func NewUserService(log log.Logger, ubiz *biz.UserUsecase) *UserService {
	return &UserService{
		log:  log,
		ubiz: ubiz,
	}
}

func pbUserInfoToBizUser(info *userpb.UserBaseInfo) *biz.User {
	return &biz.User{
		Uid:      info.Uid,
		Account:  info.Account,
		PassWD:   info.Passwd,
		PhoneNum: info.PhoneNum,
		Name:     info.Name,
		Status:   int(info.Status),
	}
}

func bizUserToPbUser(bu *biz.User) (info *userpb.UserBaseInfo) {
	return &userpb.UserBaseInfo{
		Uid:      bu.Uid,
		Account:  bu.Account,
		PhoneNum: info.PhoneNum,
		Name:     info.Name,
		Status:   userpb.UserStatus(bu.Status),
	}
}

// LoginUser implements v1.UserServer.
func (us *UserService) LoginUser(ctx context.Context, req *userpb.LoginUserReq) (resp *userpb.LoginUserResp, err error) {
	bu, e := us.ubiz.UserLogin(ctx, req.Account, req.Passwd)
	if e != nil {
		return nil, e
	}
	return &userpb.LoginUserResp{Uid: bu.Uid, Name: bu.Name}, nil
}

func (us *UserService) CreateUser(c context.Context, req *userpb.CreateUserRequest) (resp *userpb.CreateUserReply, err error) {
	if req.Info == nil {
		err = userpb.ErrorInvaildParam("info is nil")
		return
	}
	req.Info.Status = userpb.UserStatus_NOT_ACTIVE
	bu := pbUserInfoToBizUser(req.Info)
	if bu, err = us.ubiz.Create(c, bu); err != nil {
		return
	}
	req.Info.Uid = bu.Uid
	return &userpb.CreateUserReply{Info: req.Info}, nil
}

func (us *UserService) UpdateUser(c context.Context, req *userpb.UpdateUserRequest) (resp *userpb.UpdateUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (us *UserService) DeleteUser(c context.Context, req *userpb.DeleteUserRequest) (resp *userpb.DeleteUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (us *UserService) GetUser(c context.Context, req *userpb.GetUserRequest) (resp *userpb.GetUserReply, err error) {
	panic("not implemented") // TODO: Implement
}

func (us *UserService) ListUser(c context.Context, req *userpb.ListUserRequest) (resp *userpb.ListUserReply, err error) {
	panic("not implemented") // TODO: Implement
}
