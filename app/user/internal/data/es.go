package data

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

type UserESDO struct {
	Name string `json:"name"`
	Uid  int64  `json:"uid"`
}

func userIndexName() string {
	return "xshop_user"
}

func (ur *userRepo) initESIndex(ctx context.Context) error {
	yes, err := ur.data.esCli.IndexExists(userIndexName()).Do(ctx)
	if err != nil {
		log.Errorf("initESIndex IndexExists indexName:%s, err:%v", userIndexName(), err)
		return err
	}
	if !yes {
		if _, err = ur.data.esCli.CreateIndex(userIndexName()).BodyString(userIndexMapping).Do(ctx); err != nil {
			log.Errorf("initESIndex CreateIndex indexName:%s, err:%v", userIndexName(), err)
			return err
		}
	}
	return nil
}

var userIndexMapping = `
{
  "mappings" : {
	"properties" : {
		"name" : {
		  "type" : "text",
		  "analyzer": "standard"
		},
		"uid" : {
			"type" : "long"
			}
		}		  
	}
  }
}  
`

func (ud *UserDO) AfterCreate(tx *gorm.DB) error {
	esModel := UserESDO{
		Name: ud.Name,
		Uid:  ud.Uid,
	}

	_, err := esCli.Index().Index(userIndexName()).BodyJson(esModel).Id(strconv.FormatInt(ud.Uid, 10)).Do(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (ud *UserDO) AfterUpdate(tx *gorm.DB) (err error) {
	if len(ud.Name) <= 0 {
		return
	}

	esModel := UserESDO{
		Uid:  ud.Uid,
		Name: ud.Name,
	}
	_, err = esCli.Update().Index(userIndexName()).
		Doc(esModel).Id(strconv.FormatInt(ud.Uid, 10)).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (ud *UserDO) AfterDelete(tx *gorm.DB) (err error) {
	_, err = esCli.Delete().Index(userIndexName()).Id(strconv.FormatInt(ud.Uid, 10)).Do(context.Background())
	if err != nil {
		esErr := err.(*elastic.Error)
		if esErr.Status == http.StatusNotFound {
			return nil
		}
		return err
	}
	return nil
}
