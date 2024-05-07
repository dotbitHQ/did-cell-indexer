package dao

import (
	"did-cell-indexer/tables"
	"gorm.io/gorm"
)

func (d *DbDao) SearchRecordsByAccount(accountId string) (list []tables.TableRecordsInfo, err error) {
	err = d.parserDb.Where(" account_id=? ", accountId).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}
