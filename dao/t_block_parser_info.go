package dao

import (
	"did-cell-indexer/tables"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *DbDao) FindBlockInfo(parserType tables.ParserType) (block tables.TableBlockParserInfo, err error) {
	err = d.db.Where("parser_type=?", parserType).
		Order("block_number DESC").Limit(1).Find(&block).Error
	return
}

func (d *DbDao) CreateBlockInfo(parserType tables.ParserType, blockNumber uint64, blockHash, parentHash string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Insert{
			Modifier: "IGNORE",
		}).Create(&tables.TableBlockParserInfo{
			ParserType:  parserType,
			BlockNumber: blockNumber,
			BlockHash:   blockHash,
			ParentHash:  parentHash,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(tables.TableBlockParserInfo{}).
			Where("parser_type=? AND block_number=?", parserType, blockNumber).
			Updates(map[string]interface{}{
				"block_hash":  blockHash,
				"parent_hash": parentHash,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) DeleteBlockInfo(parserType tables.ParserType, blockNumber uint64) error {
	return d.db.Where("parser_type=? AND block_number < ?", parserType, blockNumber).
		Delete(&tables.TableBlockParserInfo{}).Error
}

func (d *DbDao) FindBlockInfoByBlockNumber(parserType tables.ParserType, blockNumber uint64) (block tables.TableBlockParserInfo, err error) {
	err = d.db.Where("parser_type=? AND block_number=?", parserType, blockNumber).
		Order("block_number DESC").Limit(1).Find(&block).Error
	return
}
