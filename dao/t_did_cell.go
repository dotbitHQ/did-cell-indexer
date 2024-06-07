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

func (d *DbDao) QueryDidCell(args, keyword string, limit, offset int, didType tables.DidCellStatus) (didList []tables.TableDidCellInfo, err error) {
	recycleAt := tables.GetDidCellRecycleExpiredAt()
	nowAt := time.Now().Unix()
	db := d.db.Where("args=?", args)
	switch didType {
	case tables.DidCellStatusDefault:
	case tables.DidCellStatusNormal:
		db = db.Where("expired_at>=?", nowAt)
	case tables.DidCellStatusExpired:
		db = db.Where("expired_at<=? AND expired_at>=", nowAt, recycleAt)
	case tables.DidCellStatusRecycle:
		db = db.Where("expired_at<=?", recycleAt)
	}

	if keyword != "" {
		db = db.Where("account LIKE ?", "%"+keyword+"%")
	}
	err = db.Limit(limit).Offset(offset).Find(&didList).Error
	return
}

func (d *DbDao) QueryDidCellTotal(args, keyword string, didType tables.DidCellStatus) (count int64, err error) {
	recycleAt := tables.GetDidCellRecycleExpiredAt()
	nowAt := time.Now().Unix()
	db := d.db.Model(tables.TableDidCellInfo{}).Where("args=?", args)
	switch didType {
	case tables.DidCellStatusDefault:
	case tables.DidCellStatusNormal:
		db = db.Where("expired_at>=?", nowAt)
	case tables.DidCellStatusExpired:
		db = db.Where("expired_at<=? AND expired_at>=", nowAt, recycleAt)
	case tables.DidCellStatusRecycle:
		db = db.Where("expired_at<=?", recycleAt)
	}

	if keyword != "" {
		db = db.Where("account LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&count).Error
	return
}

func (d *DbDao) GetAccountInfoByAccountId(accountId string) (acc tables.TableDidCellInfo, err error) {
	err = d.db.Where(" account_id= ? ", accountId).
		Order("expired_at DESC").Limit(1).Find(&acc).Error
	return
}

func (d *DbDao) GetAccountInfoByOutpoint(outpoint string) (acc tables.TableDidCellInfo, err error) {
	err = d.db.Where(" outpoint= ? ", outpoint).Find(&acc).Error
	return
}
