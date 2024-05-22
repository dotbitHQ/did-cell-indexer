package dao

import (
	"did-cell-indexer/tables"
	"gorm.io/gorm"
	"time"
)

func (d *DbDao) AccountUpgrade(didCellInfo tables.TableDidCellInfo) error {
	return d.db.Create(didCellInfo).Error
}

func (d *DbDao) CreateDidCellRecordsInfos(outpoint string, didCellInfo tables.TableDidCellInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id = ?", didCellInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}
		if err := tx.Select("outpoint", "block_number").
			Where("outpoint = ?", outpoint).
			Updates(didCellInfo).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) EditDidCellOwner(outpoint string, didCellInfo tables.TableDidCellInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("outpoint", "block_number", "args", "lock_code_hash").
			Where("outpoint = ?", outpoint).
			Updates(didCellInfo).Error; err != nil {
			return err
		}
		return nil

	})
}

func (d *DbDao) DidCellRecycle(outpoint string, accountId string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id=?", accountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("outpoint = ? ", outpoint).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) DidCellRenew(oldDidCellOutpoint string, didCellInfo tables.TableDidCellInfo) error {
	return d.db.Select("outpoint", "expired_at", "block_number").Where("outpoint = ?", oldDidCellOutpoint).Updates(didCellInfo).Error
}

func (d *DbDao) QueryDidCell(args string, didType tables.DidCellStatus) (didList []tables.TableDidCellInfo, err error) {
	sql := d.db.Where(" args= ?", args)
	if didType == tables.DidCellStatusNormal {
		sql.Where("expired_at > ?", time.Now().Unix())
	} else if didType == tables.DidCellStatusExpired {
		sql.Where("expired_at <= ?", time.Now().Unix())
	}
	err = sql.Find(&didList).Error
	return
}

func (d *DbDao) GetAccountInfoByAccountId(accountId string) (acc tables.TableDidCellInfo, err error) {
	err = d.db.Where(" account_id= ? ", accountId).Find(&acc).Error
	return
}

func (d *DbDao) GetAccountInfoByOutpoint(outpoint string) (acc tables.TableDidCellInfo, err error) {
	err = d.db.Where(" outpoint= ? ", outpoint).Find(&acc).Error
	return
}
