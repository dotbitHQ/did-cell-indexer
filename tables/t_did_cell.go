package tables

import (
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

const (
	TableNameDidCellInfo = "t_did_cell_info"
)

type TableDidCellInfo struct {
	Id           uint64    `json:"id" gorm:"column:id;primaryKey;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	BlockNumber  uint64    `json:"block_number" gorm:"column:block_number;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	Outpoint     string    `json:"outpoint" gorm:"column:outpoint;uniqueIndex:uk_op;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	AccountId    string    `json:"account_id" gorm:"column:account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Account      string    `json:"account" gorm:"column:account;index:k_account;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Args         string    `json:"args" gorm:"column:args;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	LockCodeHash string    `json:"lock_code_hash" gorm:"column:lock_code_hash;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' "`
	ExpiredAt    uint64    `json:"expired_at" gorm:"column:expired_at;index:k_expired_at;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

func (t *TableDidCellInfo) TableName() string {
	return TableNameDidCellInfo
}

func GetDidCellRecycleExpiredAt() uint64 {
	return uint64(time.Now().Add(-time.Hour * 24 * 30 * 3).Unix())
}

type DidCellStatus int

const (
	DidCellStatusDefault DidCellStatus = 0
	DidCellStatusNormal  DidCellStatus = 1
	DidCellStatusExpired DidCellStatus = 2
	DidCellStatusRecycle DidCellStatus = 3
)

func (t *TableDidCellInfo) IsExpired() bool {
	if int64(t.ExpiredAt) <= time.Now().Unix() {
		return true
	}
	return false
}

func (t *TableDidCellInfo) GetOutpoint() *types.OutPoint {
	return common.String2OutPointStruct(t.Outpoint)
}
