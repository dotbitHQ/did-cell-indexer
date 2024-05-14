package dao

import (
	"did-cell-indexer/config"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"gorm.io/gorm"
)

type DbDao struct {
	db *gorm.DB
}

func NewGormDB(dbMysql config.DbMysql) (*DbDao, error) {
	fmt.Println(11111, dbMysql.Addr)
	db, err := http_api.NewGormDB(dbMysql.Addr, dbMysql.User, dbMysql.Password, dbMysql.DbName, 100, 100)
	if err != nil {
		return nil, fmt.Errorf("http_api.NewGormDB err: %s", err.Error())
	}

	if err = db.AutoMigrate(
		tables.TableBlockParserInfo{},
		tables.TableDidCellInfo{},
		tables.TableRecordsInfo{},
	); err != nil {
		return nil, err
	}

	return &DbDao{db: db}, nil
}
