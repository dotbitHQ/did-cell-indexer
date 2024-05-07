package block_parser

import (
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"

	"github.com/dotbitHQ/das-lib/witness"

	"strconv"
)

func (b *BlockParser) ActionAccountUpgrade(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version account cross chain tx")
		return
	}
	log.Info("ActionAccountCrossChain:", req.BlockNumber, req.TxHash, req.Action)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
	if err != nil {
		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
		return
	}
	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[didEntity.Target.Index].Lock.Args)

	didCellInfo := tables.TableDidCellInfo{
		BlockNumber: req.BlockNumber,
		Outpoint:    common.OutPoint2String(req.TxHash, 0),
		AccountId:   builder.AccountId,
		Args:        didCellArgs,
	}

	if err = b.DbDao.AccountUpgrade(didCellInfo); err != nil {
		log.Error("AccountCrossChain err:", err.Error(), req.TxHash, req.BlockNumber)
		resp.Err = fmt.Errorf("AccountCrossChain err: %s ", err.Error())
		return
	}
	return
}

// edit record
func (b *BlockParser) ActionEditDidCellRecords(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version account cross chain tx")
		return
	}
	log.Info("ActionAccountCrossChain:", req.BlockNumber, req.TxHash, req.Action)

	txDidEntity, err := witness.TxToDidEntity(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
		return
	}

	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(req.Tx.OutputsData[txDidEntity.Outputs[0].Target.Index]); err != nil {
		resp.Err = fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
		return
	}

	var recordsInfos []tables.TableRecordsInfo
	account := didCellData.Account
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	recordList := txDidEntity.Outputs[0].DidCellWitnessDataV0.Records
	for _, v := range recordList {
		recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
			AccountId: accountId,
			Account:   account,
			Key:       v.Key,
			Type:      v.Type,
			Label:     v.Label,
			Value:     v.Value,
			Ttl:       strconv.FormatUint(uint64(v.TTL), 10),
		})
	}
	log.Info("ActionEditDidRecords:", account)

	oldDidCellOutpoint := common.OutPointStruct2String(req.Tx.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
	var didCellInfo tables.TableDidCellInfo
	didCellInfo.AccountId = accountId
	didCellInfo.BlockNumber = req.BlockNumber
	didCellInfo.Outpoint = common.OutPoint2String(req.Tx.Hash.Hex(), uint(txDidEntity.Outputs[0].Target.Index))
	if err := b.DbDao.CreateDidCellRecordsInfos(oldDidCellOutpoint, didCellInfo, recordsInfos); err != nil {
		log.Error("CreateDidCellRecordsInfos err:", err.Error())
		resp.Err = fmt.Errorf("CreateDidCellRecordsInfos err: %s", err.Error())
	}

	return
}

func (b *BlockParser) ActionEditDidCellOwner(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version account cross chain tx")
		return
	}
	log.Info("ActionAccountCrossChain:", req.BlockNumber, req.TxHash, req.Action)
	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
	if err != nil {
		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
		return
	}
	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(req.Tx.OutputsData[didEntity.Target.Index]); err != nil {
		resp.Err = fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
		return
	}
	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[didEntity.Target.Index].Lock.Args)
	account := didCellData.Account
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	didCellInfo := tables.TableDidCellInfo{
		BlockNumber: req.BlockNumber,
		Outpoint:    common.OutPoint2String(req.TxHash, 0),
		AccountId:   accountId,
		Args:        didCellArgs,
	}

	oldOutpoint := common.OutPointStruct2String(req.Tx.Inputs[0].PreviousOutput)
	if err := b.DbDao.EditDidCellOwner(oldOutpoint, didCellInfo); err != nil {
		log.Error("EditDidCellOwner err:", err.Error())
		resp.Err = fmt.Errorf("EditDidCellOwner err: %s", err.Error())
	}
	return
}

func (b *BlockParser) ActionRenewAccount(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version renew account tx")
		return
	}
	log.Info("ActionRenewAccount:", req.BlockNumber, req.TxHash)

	builder, err := witness.AccountCellDataBuilderFromTx(req.Tx, common.DataTypeNew)
	if err != nil {
		resp.Err = fmt.Errorf("AccountCellDataBuilderFromTx err: %s", err.Error())
		return
	}

	log.Info("ActionRenewAccount:", builder.Account, builder.ExpiredAt)

	//renew did cell
	var didCellInfo tables.TableDidCellInfo
	var oldDidCellOutpoint string

	var didCellData witness.DidCellData
	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
	if err != nil {
		resp.Err = fmt.Errorf("witness.TxToOneDidEntity err: %s", err.Error())
		return
	}
	if err := didCellData.BysToObj(req.Tx.OutputsData[didEntity.Target.Index]); err != nil {
		resp.Err = fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
		return
	}
	oldDidCellOutpoint = common.OutPointStruct2String(req.Tx.Inputs[didEntity.Target.Index].PreviousOutput)
	didCellInfo.Outpoint = common.OutPoint2String(req.Tx.Hash.Hex(), uint(didEntity.Target.Index))
	didCellInfo.ExpiredAt = didCellData.ExpireAt
	didCellInfo.BlockNumber = req.BlockNumber

	if err := b.DbDao.DidCellRenew(oldDidCellOutpoint, didCellInfo); err != nil {
		log.Error("RenewAccount err:", err.Error())
		resp.Err = fmt.Errorf("RenewAccount err: %s", err.Error())
	}
	return
}

func (b *BlockParser) ActionDidCellRecycle(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
		return
	} else if !isCV {
		log.Warn("not current version account cross chain tx")
		return
	}
	log.Info("ActionAccountCrossChain:", req.BlockNumber, req.TxHash, req.Action)
	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
	if err != nil {
		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
		return
	}
	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(req.Tx.OutputsData[didEntity.Target.Index]); err != nil {
		resp.Err = fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
		return
	}
	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[didEntity.Target.Index].Lock.Args)
	account := didCellData.Account
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	var didCellInfo tables.TableDidCellInfo
	didCellInfo.Args = didCellArgs
	didCellInfo.AccountId = accountId
	oldOutpoint := common.OutPointStruct2String(req.Tx.Inputs[0].PreviousOutput)
	if err := b.DbDao.DidCellRecycle(oldOutpoint); err != nil {
		log.Error("DidCellRecycle err:", err.Error())
		resp.Err = fmt.Errorf("DidCellRecycle err: %s", err.Error())
	}
	return

}
