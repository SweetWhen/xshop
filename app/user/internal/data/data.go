package data

import (
	"fmt"
	"realworld/app/user/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserDO struct {
	Uid       int64  `gorm:"primaryKey"`
	Account   string `gorm:"index:idx_account,unique"`
	PassWD    string
	Name      string
	PhoneNum  string
	Status    int
	CreatedAt time.Time
	UpatedAt  time.Time
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewUserRepo, NewData)

// Data .
type Data struct {
	userDB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User, c.Database.Passwd, c.Database.Host, c.Database.Port, c.Database.DbName)
	log.Infof("userdata dsn: %s", dsn)
	d, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, nil, err
	}
	d2, _ := d.DB()
	d2.SetMaxIdleConns(10)
	d2.SetMaxOpenConns(100)
	d2.SetConnMaxIdleTime(time.Minute * 30)

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:         conf.Redis.Addr,
	// 	Password:     conf.Redis.Password,
	// 	DB:           int(conf.Redis.Db),
	// 	DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
	// 	WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
	// 	ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	// })
	// rdb.AddHook(redisotel.TracingHook{})

	resData := &Data{
		userDB: d,
		// rdb: rdb,
	}
	if err := resData.userDB.AutoMigrate(&UserDO{}); err != nil {
		return nil, nil, err
	}

	return resData, func() {
		log.Info("message", "closing the data resources")
		if err := d2.Close(); err != nil {
			log.Error(err)
		}
		// if err := d.rdb.Close(); err != nil {
		// 	log.Error(err)
		// }
	}, nil
}
