package data

import (
	"context"
	mgrpb "realworld/api/mgr/v1"
	userpb "realworld/api/user/v1"
	"realworld/app/mgr/internal/biz"

	"google.golang.org/grpc"
)

func NewUserProxy(d *Data) biz.IUser {
	return &UserProxy{d: d}
}

func userToSvrUser(in *biz.User) *userpb.UserBaseInfo {
	return &userpb.UserBaseInfo{
		Uid:      in.Uid,
		Account:  in.Account,
		Passwd:   in.PassWD,
		PhoneNum: in.PhoneNum,
		Name:     in.Name,
		Status:   userpb.UserStatus(in.Status),
	}
}

func userSvrToUser(info *userpb.UserBaseInfo) *biz.User {
	return &biz.User{
		Uid:      info.Uid,
		Account:  info.Account,
		PassWD:   info.Passwd,
		PhoneNum: info.PhoneNum,
		Name:     info.Name,
		Status:   int(info.Status),
	}
}

var _ biz.IUser = &UserProxy{}

type UserProxy struct {
	d *Data
}

// DeleteUser implements biz.IUser.
func (up *UserProxy) UserLogin(ctx context.Context, ac, pw string) (*biz.User, error) {
	resp, err := up.d.userCli.LoginUser(ctx, &userpb.LoginUserReq{Account: ac, Passwd: pw})
	if err != nil {
		return nil, err
	}

	return &biz.User{Uid: resp.Uid, Name: resp.Name}, nil
}

// CreateUser implements biz.IUser.
func (up *UserProxy) CreateUser(ctx context.Context, in *biz.User, opts ...grpc.CallOption) (*biz.User, error) {
	req := userpb.CreateUserRequest{Info: userToSvrUser(in)}
	resp, err := up.d.userCli.CreateUser(ctx, &req)
	if err != nil {
		return nil, err
	}
	if resp.Info == nil {
		return nil, mgrpb.ErrorUsersvrBadResp("userSvr.CreateUser resp.Info is nil")
	}
	return userSvrToUser(resp.Info), nil
}

// DeleteUser implements biz.IUser.
func (up *UserProxy) DeleteUser(ctx context.Context, in *biz.User, opts ...grpc.CallOption) (*biz.User, error) {
	panic("unimplemented")
}

// GetUser implements biz.IUser.
func (up *UserProxy) GetUser(ctx context.Context, in *biz.User, opts ...grpc.CallOption) (*biz.User, error) {
	panic("unimplemented")
}

// ListUser implements biz.IUser.
func (up *UserProxy) ListUser(ctx context.Context, in *biz.User, opts ...grpc.CallOption) (*biz.User, error) {
	panic("unimplemented")
}

// UpdateUser implements biz.IUser.
func (up *UserProxy) UpdateUser(ctx context.Context, in *biz.User, opts ...grpc.CallOption) (*biz.User, error) {
	panic("unimplemented")
}
