package data

import (
	"context"
	"errors"

	userpb "realworld/api/user/v1"
	"realworld/app/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(d *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: d,
		log:  log.NewHelper(logger),
	}
}

func (ur *userRepo) Create(ctx context.Context, bu *biz.User) (*biz.User, error) {
	ud := bizUserToDOUser(bu)
	result := ur.data.userDB.WithContext(ctx).Create(ud)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		// duplicate account
		return nil, userpb.ErrorUserExisted("account:%s exited", bu.Account)
	} else if result.Error != nil {
		return nil, result.Error
	}
	bu.Uid = ud.Uid
	return bu, nil
}

func bizUserToDOUser(bu *biz.User) *UserDO {
	return &UserDO{
		Uid:      bu.Uid,
		Account:  bu.Account,
		PassWD:   bu.PassWD,
		Name:     bu.Name,
		PhoneNum: bu.PhoneNum,
		Status:   bu.Status,
	}
}

func doUserToBizUser(ud *UserDO) *biz.User {
	return &biz.User{
		Uid:      ud.Uid,
		Account:  ud.Account,
		PassWD:   ud.PassWD,
		Name:     ud.Name,
		PhoneNum: ud.PhoneNum,
		Status:   ud.Status,
	}
}

func (ur *userRepo) Update(ctx context.Context, bu *biz.User) error {
	ud := bizUserToDOUser(bu)
	result := ur.data.userDB.WithContext(ctx).Model(&UserDO{}).Updates(ud)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *userRepo) Get(ctx context.Context, ac string) (*biz.User, error) {
	ud := UserDO{}
	result := ur.data.userDB.WithContext(ctx).Where("account = ?", ac).First(&ud)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, userpb.ErrorUserNotFound("account: %s", ac)
	}

	return doUserToBizUser(&ud), nil
}

func (ur *userRepo) ListUser(ctx context.Context, startId int64, cnt int64) (bus []*biz.User, nextStartId int64, err error) {
	users := make([]*UserDO, 0, int(cnt))
	result := ur.data.userDB.WithContext(ctx).Where("uid > ?", startId).Order("uid asc").Limit(int(cnt)).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	bus = make([]*biz.User, 0, len(users))
	for _, u := range users {
		bus = append(bus, doUserToBizUser(u))
	}
	if len(users) == int(cnt) {
		nextStartId = users[len(users)-1].Uid
	}
	return
}
