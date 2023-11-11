package biz

import (
	"context"

	userpb "realworld/api/user/v1"

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

// UserRepo is a User repo.
type UserRepo interface {
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) error
	Get(context.Context, string) (*User, error)
	ListUser(ctx context.Context, startId int64, cnt int64) (bus []*User, nextStartId int64, err error)
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

var _ UserRepo = &UserUsecase{}

// UserUsecase is a Greeter usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// UserLogin implements UserRepo.
func (uc *UserUsecase) UserLogin(ctx context.Context, account string, pw string) (u *User, err error) {
	if u, err = uc.repo.Get(ctx, account); err != nil {
		return
	}
	if u.PassWD != pw || u.Status != int(userpb.UserStatus_ACTIVE) {
		err = userpb.ErrorWrongPasswd("account: %s, staus:%d", account, u.Status)
		return
	}
	u.PassWD = ""
	return
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) Create(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Account)
	return uc.repo.Create(ctx, g)
}

func (uc *UserUsecase) Update(c context.Context, g *User) error {
	return uc.repo.Update(c, g)
}

func (uc *UserUsecase) Get(c context.Context, ac string) (*User, error) {
	return uc.repo.Get(c, ac)
}

func (uc *UserUsecase) ListUser(ctx context.Context, startId int64, cnt int64) (bus []*User, nextStartId int64, err error) {
	return uc.repo.ListUser(ctx, startId, cnt)
}
