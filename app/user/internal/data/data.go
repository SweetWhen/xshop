package data

import (
	"fmt"
	fmtLog "log"
	"os"
	"time"

	"realworld/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

type UserDO struct {
	Uid       int64     `gorm:"primaryKey"`
	Account   string    `gorm:"uniqueIndex:idx_uk_account;type:varchar(128) not null;defualt:'';comment:'用户账号'"`
	PassWD    string    `gorm:"type:varchar(256) not null;default:'';comment:'用户密码'"`
	Name      string    `gorm:"type:varchar(32) not null;default:''"`
	PhoneNum  string    `gorm:"type:varchar(16) not null; defualt:''"`
	Status    int       `gorm:"type:smallint;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpatedAt  time.Time `gorm:"autoUpdateTime:nano"`
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewUserRepo, NewData)

// Data .
type Data struct {
	userDB *gorm.DB
	esCli  *elastic.Client
}

var esCli *elastic.Client

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// return &Data{}, func() {}, nil
	log := log.NewHelper(logger)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User, c.Database.Passwd, c.Database.Host, c.Database.Port, c.Database.DbName)
	log.Infof("userdata dsn: %s", dsn)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormLog.Default.LogMode(gormLog.Info),
	})
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
	resData.esCli, err = InitEs(c.EsInfo)
	if err != nil {
		panic(err)
	}
	esCli = resData.esCli
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

func InitEs(esCfg *conf.Data_ES) (cli *elastic.Client, err error) {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", esCfg.Host, esCfg.Port)
	cli, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(fmtLog.New(os.Stdout, "xshop", fmtLog.Flags())))
	if err != nil {
		return
	}
	return
}
