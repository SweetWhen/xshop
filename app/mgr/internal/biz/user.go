package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
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

// GreeterUsecase is a Greeter usecase.
type MgrUserUsecase struct {
	iu  IUser
	log *log.Helper
}

type IUser interface {
	UserLogin(ctx context.Context, ac, pw string) (*User, error)
	CreateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	UpdateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	DeleteUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	GetUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	ListUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
}

// NewMgrUserUsecase new a Greeter usecase.
func NewMgrUserUsecase(iu IUser, logger log.Logger) *MgrUserUsecase {
	return &MgrUserUsecase{iu: iu, log: log.NewHelper(logger)}
}

func (mu *MgrUserUsecase) UserLogin(ctx context.Context, account string, pw string) (u *User, err error) {
	resp, err2 := mu.iu.UserLogin(ctx, account, pw)
	if err2 != nil {
		err = err2
		return
	}
	return resp, nil
}

func (mu *MgrUserUsecase) CreateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (u *User, err error) {
	resp, err2 := mu.iu.CreateUser(ctx, in)
	if err2 != nil {
		err = err2
		return
	}
	return resp, nil
}

func (mu *MgrUserUsecase) UpdateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (u *User, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserUsecase) DeleteUser(ctx context.Context, in *User, opts ...grpc.CallOption) (u *User, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserUsecase) GetUser(ctx context.Context, in *User, opts ...grpc.CallOption) (u *User, err error) {
	panic("not implemented") // TODO: Implement
}

func (mu *MgrUserUsecase) ListUser(ctx context.Context, in *User, opts ...grpc.CallOption) (u *User, err error) {
	panic("not implemented") // TODO: Implement
}
