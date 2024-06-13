package dao

import (
	"did-cell-indexer/config"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"gorm.io/gorm"
)

type DbDao struct {
	db *gorm.DB
}

func NewGormDB(dbMysql config.DbMysql) (*DbDao, error) {
	log := logger.NewLoggerDefault("dao", logger.LevelDebug, nil)
	db, err := http_api.NewGormDBWithLog(dbMysql.Addr, dbMysql.User, dbMysql.Password, dbMysql.DbName, 100, 100, log)
	if err != nil {
		return nil, fmt.Errorf("http_api.NewGormDB err: %s", err.Error())
	}

	if err = db.AutoMigrate(
		tables.TableBlockParserInfo{},
		tables.TableDidCellInfo{},
		tables.TableRecordsInfo{},
		tables.TableTxInfo{},
	); err != nil {
		return nil, err
	}

	return &DbDao{db: db}, nil
}
