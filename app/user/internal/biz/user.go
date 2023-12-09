package biz

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"

	userpb "realworld/api/user/v1"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-kratos/kratos/v2/log"
)

// User is a User model.
type User struct {
	Uid      int64
	Account  string
	PassWD   string
	PhoneNum string
	Name     string
	Status   int
}

type UserUpdate struct {
	Account  string
	PassWD   *string
	PhoneNum *string
	Name     *string
}

// UserRepo is a User repo.
type UserRepo interface {
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *UserUpdate) error
	Get(context.Context, string) (*User, error)
	Delete(context.Context, string) error
	ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*User, nextStartId int64, err error)
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		log:     log.NewHelper(logger),
		rsaImpl: NewRSAImpl(""),
	}
}

var _ UserRepo = &UserUsecase{}

// UserUsecase is a Greeter usecase.
type UserUsecase struct {
	repo    UserRepo
	rsaImpl *RSAImpl
	log     *log.Helper
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

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(u.PassWD, "$")
	if len(passwordInfo) != 2 {
		uc.log.WithContext(ctx).Errorf("CreateUser: %s get a pwInfo:V", account, passwordInfo)
		err = userpb.ErrorWrongPasswd("account: %s, len(passwordInfo):%d", account, len(passwordInfo))
		return
	}
	check := password.Verify(pw, passwordInfo[0], passwordInfo[1], options)
	if !check {
		uc.log.WithContext(ctx).Debugf("CreateUser: %s  check failed", account)
		err = userpb.ErrorWrongPasswd("account: %s, check failed", account)
		return
	}
	u.PassWD = ""
	return
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) Create(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Account)
	//todo:  rsa解密用户上传的密码

	//密码加盐存储到db
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(g.PassWD, options)
	g.PassWD = fmt.Sprintf("%s$%s", salt, encodedPwd)

	u, err := uc.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}
	u.PassWD = ""
	return u, nil
}

func (uc *UserUsecase) Update(c context.Context, g *UserUpdate) error {
	return uc.repo.Update(c, g)
}

func (uc *UserUsecase) Get(c context.Context, ac string) (*User, error) {
	return uc.repo.Get(c, ac)
}

func (uc *UserUsecase) Delete(c context.Context, ac string) error {
	return uc.repo.Delete(c, ac)
}

func (uc *UserUsecase) ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*User, nextStartId int64, err error) {
	return uc.repo.ListUser(ctx, startId, cnt, int(status))
}
