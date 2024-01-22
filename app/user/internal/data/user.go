package data

import (
	"context"
	"encoding/json"
	"errors"

	userpb "realworld/api/user/v1"
	"realworld/app/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(d *Data, logger log.Logger) biz.UserRepo {
	ur := &userRepo{
		data: d,
		log:  log.NewHelper(logger),
	}
	ctx := context.Background()
	if err := ur.initESIndex(ctx); err != nil {
		panic(err)
	}
	return ur
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
	//write to es

	return bu, nil
}

func doUserToESUser(ud *UserDO) *UserESDO {
	return &UserESDO{
		Uid:  ud.Uid,
		Name: ud.Name,
	}
}

func esUserToDOUser(bu *UserESDO) *UserDO {
	return &UserDO{
		Uid:  bu.Uid,
		Name: bu.Name,
	}
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

func (ur *userRepo) Update(ctx context.Context, bu *biz.UserUpdate) error {
	ud := &UserDO{Uid: bu.Uid}
	us := make(map[string]any)
	if bu.Name != nil {
		us["name"] = *bu.Name
		ud.Name = *bu.Name
	}
	if bu.PassWD != nil {
		us["pass_wd"] = *bu.PassWD
	}
	if bu.PhoneNum != nil {
		us["phone_num"] = *bu.PhoneNum
	}
	if len(us) == 0 {
		return userpb.ErrorInvaildParam("nothing to update")
	}
	result := ur.data.userDB.WithContext(ctx).Model(ud).Where("uid = ?", bu.Uid).Updates(us)
	if result.RowsAffected == 0 {
		return userpb.ErrorUserNotFound("uid: %d", bu.Uid)
	}
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

func (ur *userRepo) SearchUser(ctx context.Context, nameKey string, cnt int32) ([]*biz.User, int32, error) {
	q := elastic.NewMatchQuery("name", nameKey)
	hl := elastic.NewHighlight()
	hl.Fields(elastic.NewHighlighterField("name")).PostTags("</em>").PreTags("<em>")
	switch {
	case cnt > 100:
		cnt = 100
	case cnt <= 0:
		cnt = 10
	}
	result, err := ur.data.esCli.Search().Index(userIndexName()).Query(q).Highlight(hl).Size(int(cnt)).Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	users := make([]*biz.User, 0, int(cnt))
	total := int32(result.Hits.TotalHits.Value)
	for _, value := range result.Hits.Hits {
		esUser := UserESDO{}
		_ = json.Unmarshal(value.Source, &esUser)
		u := &biz.User{
			Uid:  esUser.Uid,
			Name: esUser.Name,
		}
		u.Highlight = make(map[string][]string, len(value.Highlight))
		for k, v := range value.Highlight {
			u.Highlight[k] = v
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (ur *userRepo) Delete(ctx context.Context, ac string, hard int32) error {
	var result *gorm.DB
	if hard == 1 {
		result = ur.data.userDB.WithContext(ctx).Where("account = ?", ac).Delete(&UserDO{})
	} else {
		result = ur.data.userDB.WithContext(ctx).Model(&UserDO{}).Where("account = ?", ac).Update("status", userpb.UserStatus_NOT_ACTIVE)
		if result.RowsAffected == 0 {
			return userpb.ErrorUserNotFound("account: %s", ac)
		}
	}
	return result.Error
}

func (ur *userRepo) ListUser(ctx context.Context, startId int64, cnt int64, status int) (bus []*biz.User, nextStartId int64, err error) {
	users := make([]*UserDO, 0, int(cnt))
	tx := ur.data.userDB.WithContext(ctx).Where("uid > ?", startId)
	if status != 0 {
		tx.Where("status = ?", status)
	}
	result := tx.Order("uid asc").Limit(int(cnt)).Find(&users)
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
