package tables

import (
	"github.com/dotbitHQ/das-lib/common"
	"time"
)

type TableTxInfo struct {
	Id             uint64               `json:"id" gorm:"column:id; primaryKey; type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	Outpoint       string               `json:"outpoint" gorm:"column:outpoint; uniqueIndex:uk_outpoint; type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Action         common.DidCellAction `json:"action" gorm:"column:action; uniqueIndex:uk_outpoint; index:k_action; type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	BlockNumber    uint64               `json:"block_number" gorm:"column:block_number; type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	BlockTimestamp int64                `json:"block_timestamp" gorm:"column:block_timestamp; index:k_block_timestamp; type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	AccountId      string               `json:"account_id" gorm:"account_id; index:k_account_id; type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Account        string               `json:"account" gorm:"column:account; index:k_account; type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Args           string               `json:"args" gorm:"column:args; index:k_args; type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Address        string               `json:"address" gorm:"column:address; type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	LockCodeHash   string               `json:"lock_code_hash" gorm:"column:lock_code_hash; type:varchar(255) NOT NULL DEFAULT '' "`
	CreatedAt      time.Time            `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt      time.Time            `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

const (
	TableNameTxInfo = "t_tx_info"
)

func (t *TableTxInfo) TableName() string {
	return TableNameTxInfo
}
