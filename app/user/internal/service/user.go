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

func NewUserService(l log.Logger, ubiz *biz.UserUsecase) *UserService {
	return &UserService{
		log:  l,
		ubiz: ubiz,
	}
}
func pbUserUpdateToBizUpdate(info *userpb.UpdateUserRequest) *biz.UserUpdate {
	return &biz.UserUpdate{
		Uid:      info.Uid,
		PassWD:   info.Passwd,
		PhoneNum: info.PhoneNum,
		Name:     info.Name,
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
		PhoneNum: bu.PhoneNum,
		Name:     bu.Name,
		Status:   userpb.UserStatus(bu.Status),
	}
}

func bizHLToPb(bu *biz.User) []*userpb.SearchUserResp_HL {
	if len(bu.Highlight) <= 0 {
		return nil
	}
	r := make([]*userpb.SearchUserResp_HL, 0, len(bu.Highlight))

	for k, v := range bu.Highlight {
		hl := &userpb.SearchUserResp_HL{
			Feild:  k,
			Values: v,
		}
		r = append(r, hl)
	}

	return r
}

func (us *UserService) GetJWTSK() string {
	return us.ubiz.GetJWTSK()
}

func (us *UserService) GetJWTPK() string {
	return us.ubiz.GetJWTPK()
}

// GetLoginInfo implements v1.UserServer.
func (us *UserService) GetLoginInfo(ctx context.Context, req *userpb.LoginInfoRequest) (*userpb.LoginInfoReply, error) {
	pk, sk := us.ubiz.GetLoginInfo(ctx, req.Account)
	return &userpb.LoginInfoReply{PublicKey: pk, PrivateKey: sk}, nil
}

// SearchUser implements v1.UserServer.
func (us *UserService) SearchUser(ctx context.Context, req *userpb.SearchUserReq) (*userpb.SearchUserResp, error) {
	bizUser, total, err := us.ubiz.SearchUser(ctx, req.NameKey, req.Count)
	if err != nil {
		return nil, err
	}
	resp := &userpb.SearchUserResp{Total: total}
	for _, bu := range bizUser {
		pbUser := bizUserToPbUser(bu)
		hl := bizHLToPb(bu)
		ui := &userpb.SearchUserResp_UserInfo{Info: pbUser, Hl: hl}
		resp.Users = append(resp.Users, ui)
	}
	return resp, nil
}

// LoginUser implements v1.UserServer.
func (us *UserService) LoginUser(ctx context.Context, req *userpb.LoginUserReq) (resp *userpb.LoginUserResp, err error) {
	bu, e := us.ubiz.UserLogin(ctx, req.Account, req.Passwd)
	if e != nil {
		return nil, e
	}
	return &userpb.LoginUserResp{Uid: bu.Uid, Name: bu.Name}, nil
}

// LogoutUser implements v1.UserServer.
func (us *UserService) LogoutUser(ctx context.Context, req *userpb.LogoutUserReq) (*userpb.LogoutUserResp, error) {

	err := us.ubiz.LogoutUser(ctx, req.Account)

	return &userpb.LogoutUserResp{}, err
}

func (us *UserService) CreateUser(c context.Context, req *userpb.CreateUserRequest) (resp *userpb.CreateUserReply, err error) {
	if req.Info == nil {
		err = userpb.ErrorInvaildParam("info is nil")
		return
	}
	req.Info.Status = userpb.UserStatus_ACTIVE
	bu := pbUserInfoToBizUser(req.Info)
	if bu, err = us.ubiz.Create(c, bu); err != nil {
		return
	}
	req.Info.Uid = bu.Uid
	req.Info.Passwd = ""
	return &userpb.CreateUserReply{Info: req.Info}, nil
}

func (us *UserService) UpdateUser(c context.Context, req *userpb.UpdateUserRequest) (resp *userpb.UpdateUserReply, err error) {
	err = us.ubiz.Update(c, pbUserUpdateToBizUpdate(req))
	return &userpb.UpdateUserReply{}, err
}

func (us *UserService) DeleteUser(c context.Context, req *userpb.DeleteUserRequest) (resp *userpb.DeleteUserReply, err error) {
	err = us.ubiz.Delete(c, req.Account, req.Hard)
	return &userpb.DeleteUserReply{}, err
}

func (us *UserService) GetUser(c context.Context, req *userpb.GetUserRequest) (resp *userpb.GetUserReply, err error) {
	bu, err := us.ubiz.Get(c, req.Account)
	if err != nil {
		return nil, err
	}
	return &userpb.GetUserReply{Info: bizUserToPbUser(bu)}, nil
}

func (us *UserService) ListUser(c context.Context, req *userpb.ListUserRequest) (resp *userpb.ListUserReply, err error) {
	if req.StartId < 0 {
		req.StartId = 0
	}
	if req.Count <= 0 {
		req.Count = 200
	}
	bus, nextStartId, err := us.ubiz.ListUser(c, req.StartId, req.Count, int(req.Status))
	if err != nil {
		return nil, err
	}
	resp = &userpb.ListUserReply{
		NextStartId: nextStartId,
	}
	for _, u := range bus {
		bu := bizUserToPbUser(u)
		resp.Users = append(resp.Users, bu)
	}
	return
}
