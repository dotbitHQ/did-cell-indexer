package block_parser

import (
	"did-cell-indexer/dao"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func (b *BlockParser) registerTransactionHandle() {
	b.mapTransactionHandle = make(map[string]FuncTransactionHandle)
	b.mapTransactionHandle[common.DasActionConfig] = b.ActionConfig

	b.mapTransactionHandle[common.DasActionAccountCellUpgrade] = b.ActionAccountUpgrade  // upgrade account cell
	b.mapTransactionHandle[common.DidCellActionEditRecords] = b.ActionEditDidCellRecords // edit did cell record
	b.mapTransactionHandle[common.DidCellActionEditOwner] = b.ActionEditDidCellOwner     // edit did cell owner
	b.mapTransactionHandle[common.DidCellActionRecycle] = b.ActionDidCellRecycle         //  did cell recycle
	b.mapTransactionHandle[common.DidCellActionRenew] = b.ActionDidCellRenew
}

func isCurrentVersionTx(tx *types.Transaction, name common.DasContractName) (bool, error) {
	contract, err := core.GetDasContractInfo(name)
	if err != nil {
		return false, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	isCV := false
	for _, v := range tx.Outputs {
		if v.Type == nil {
			continue
		}
		if contract.IsSameTypeId(v.Type.CodeHash) {
			isCV = true
			break
		}
	}
	return isCV, nil
}

type FuncTransactionHandleReq struct {
	DbDao          *dao.DbDao
	Tx             *types.Transaction
	TxHash         string
	BlockNumber    uint64
	BlockTimestamp int64
	Action         common.DasAction
}

type FuncTransactionHandleResp struct {
	ActionName string
	Err        error
}

type FuncTransactionHandle func(FuncTransactionHandleReq) FuncTransactionHandleResp
