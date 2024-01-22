package biz

import (
	"context"
	"net/http"
	"time"

	userpb "realworld/api/user/v1"
	"realworld/api/user/v1/common"
	"realworld/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// User is a User model.
type User struct {
	Uid      int64
	Account  string
	PassWD   string
	PhoneNum string
	Name     string
	Status   int

	Highlight map[string][]string
}

type UserUpdate struct {
	Uid      int64
	PassWD   *string
	PhoneNum *string
	Name     *string
}

// UserRepo is a User repo.
type UserRepo interface {
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *UserUpdate) error
	Get(context.Context, string) (*User, error)
	SearchUser(ctx context.Context, nameKey string, cnt int32) ([]*User, int32, error)
	Delete(context.Context, string, int32) error
	ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*User, nextStartId int64, err error)
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(confServer *conf.Server, repo UserRepo, logger log.Logger, confData *conf.Data) *UserUsecase {
	uc := &UserUsecase{
		repo:       repo,
		confServer: confServer,
		log:        log.NewHelper(logger),
		rsaImpl:    NewRSAImpl(confData.RsaPrivate),
		pwEncode:   NewPWEncode(),
	}

	uc.log.Debugf("NewUserUsecase privateRSA:\n%s\n, \npublicRSA:\n%s\n", confData.RsaPrivate, confData.RsaPublic)

	return uc
}

var _ UserRepo = &UserUsecase{}

// UserUsecase is a Greeter usecase.
type UserUsecase struct {
	repo       UserRepo
	rsaImpl    *RSAImpl
	pwEncode   *PWEncode
	log        *log.Helper
	confServer *conf.Server
}

func (uc *UserUsecase) GetJWTPK() string {
	return uc.rsaImpl.publicKey
}

func (uc *UserUsecase) GetJWTSK() string {
	return uc.rsaImpl.privateKey
}

func (uc *UserUsecase) GetLoginInfo(ctx context.Context, account string) (pk, sk string) {
	return uc.GetJWTPK(), uc.GetJWTSK()
}

func (uc *UserUsecase) SearchUser(ctx context.Context, nameKey string, cnt int32) ([]*User, int32, error) {
	return uc.repo.SearchUser(ctx, nameKey, cnt)
}

// UserLogin implements UserRepo.
func (uc *UserUsecase) UserLogin(ctx context.Context, account string, pw string) (u *User, err error) {
	if u, err = uc.repo.Get(ctx, account); err != nil {
		return
	}
	if u.Status != int(userpb.UserStatus_ACTIVE) {
		err = userpb.ErrorWrongPasswd("account: %s, staus:%d", account, u.Status)
		return
	}

	//todo: rsa解密得到用户密码原文
	e := uc.pwEncode.Decode(pw, u.PassWD)
	if e != nil {
		uc.log.WithContext(ctx).Errorf("UserLogin: %s get a pwInfo:V", account, u)
		err = userpb.ErrorWrongPasswd("account: %s, check failed", account)
		return
	}
	pc := userpb.ClaimPayload{
		Uid:  u.Uid,
		Name: u.Name,
	}
	var sid string
	sid, err = uc.Checkin(ctx, &pc)
	if err != nil {
		return
	}
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		err = userpb.ErrorInvaildParam("transport not found")
		return
	}
	sidCookie := &http.Cookie{
		Name:    "sid",
		Value:   sid,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	}
	tr.ReplyHeader().Set("cookie", sidCookie.String())
	uc.log.Debugf("UserLogin success account:%s, set cookie: %s", account, sidCookie.String())
	u.PassWD = ""
	return
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, ac string) error {
	claims, ok := common.GetJWTClaim(ctx)
	if !ok {
		return userpb.ErrorUserNotFound("claims not found")
	}
	uc.logout(ctx, claims.Uid)

	// todo: publish event

	return nil
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) Create(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Account)
	//todo:  rsa解密用户上传的密码

	//密码加盐存储到db
	g.PassWD = uc.pwEncode.Encode(g.PassWD)
	u, err := uc.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}
	u.PassWD = ""
	return u, nil
}

func (uc *UserUsecase) Update(c context.Context, g *UserUpdate) error {
	if g.Name != nil {
		if len(*g.Name) <= 0 {
			return userpb.ErrorInvaildParam("name can not set to empty")
		}
	}
	if g.PassWD != nil {
		newPw := *g.PassWD
		if len(newPw) <= 0 {
			return userpb.ErrorInvaildParam("pass word can not set to empty")
		}
		newPw = uc.pwEncode.Encode(newPw)
		g.PassWD = &newPw
	}
	return uc.repo.Update(c, g)
}

func (uc *UserUsecase) Get(c context.Context, ac string) (*User, error) {
	claims, ok := common.GetJWTClaim(c)
	if !ok {
		claims = &common.MyClaims{}
	}
	uc.log.Debugf("UserUsecase Get account:%s, uid:%d, name:%s",
		ac, claims.Uid, claims.Name)
	return uc.repo.Get(c, ac)
}

func (uc *UserUsecase) Delete(c context.Context, ac string, hard int32) error {
	return uc.repo.Delete(c, ac, hard)
}

func (uc *UserUsecase) ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*User, nextStartId int64, err error) {
	return uc.repo.ListUser(ctx, startId, cnt, int(status))
}
