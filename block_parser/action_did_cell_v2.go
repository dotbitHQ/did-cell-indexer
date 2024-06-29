package block_parser

import (
	"bytes"
	"did-cell-indexer/config"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"strconv"
)

func (b *BlockParser) DidCellActionUpgrade(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {

	log.Info("DidCellActionUpgrade:", req.BlockNumber, req.TxHash, req.Action)
	var list []tables.TableDidCellInfo
	var listTx []tables.TableTxInfo
	var records []tables.TableRecordsInfo
	var accountIds []string

	txDidEntityWitness, err := witness.GetDidEntityFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("witness.GetDidEntityFromTx err: %s", err.Error())
		return
	}

	for k, v := range req.TxDidCellMap.Outputs {
		_, cellData, err := v.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo err: %s[%s]", err.Error(), k)
			return
		}
		account := cellData.Account
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
		accountIds = append(accountIds, accountId)

		tmp := tables.TableDidCellInfo{
			BlockNumber:  req.BlockNumber,
			Outpoint:     common.OutPointStruct2String(v.OutPoint),
			AccountId:    accountId,
			Account:      account,
			Args:         common.Bytes2Hex(v.Lock.Args),
			LockCodeHash: v.Lock.CodeHash.Hex(),
			ExpiredAt:    cellData.ExpireAt,
		}
		list = append(list, tmp)

		mode := address.Mainnet
		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		addr, err := address.ConvertScriptToAddress(mode, v.Lock)
		if err != nil {
			resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
			return
		}
		tmpTx := tables.TableTxInfo{
			Outpoint:       common.OutPointStruct2String(v.OutPoint),
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			AccountId:      accountId,
			Account:        account,
			Action:         req.Action,
			Args:           common.Bytes2Hex(v.Lock.Args),
			Address:        addr,
			LockCodeHash:   v.Lock.CodeHash.Hex(),
		}
		listTx = append(listTx, tmpTx)

		if w, ok := txDidEntityWitness.Outputs[v.Index]; ok {
			for _, r := range w.DidCellWitnessDataV0.Records {
				records = append(records, tables.TableRecordsInfo{
					AccountId: accountId,
					Account:   account,
					Key:       r.Key,
					Type:      r.Type,
					Label:     r.Label,
					Value:     r.Value,
					Ttl:       strconv.FormatUint(uint64(r.TTL), 10),
				})
			}
		}
	}

	if err := b.DbDao.AccountUpgradeList(list, listTx, records, accountIds); err != nil {
		resp.Err = fmt.Errorf("AccountUpgradeList err: %s ", err.Error())
		return
	}

	return
}

func (b *BlockParser) DidCellActionUpdate(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	log.Info("DidCellActionUpdate:", req.BlockNumber, req.TxHash, req.Action)
	if len(req.TxDidCellMap.Inputs) != len(req.TxDidCellMap.Outputs) {
		resp.Err = fmt.Errorf("len(req.TxDidCellMap.Inputs)!=len(req.TxDidCellMap.Outputs)")
		return
	}
	txDidEntityWitness, err := witness.GetDidEntityFromTx(req.Tx)
	if err != nil {
		resp.Err = fmt.Errorf("witness.GetDidEntityFromTx err: %s", err.Error())
		return
	}

	var oldOutpointList []string
	var list []tables.TableDidCellInfo
	var accountIds []string
	var records []tables.TableRecordsInfo
	var listTx []tables.TableTxInfo

	for k, v := range req.TxDidCellMap.Inputs {
		_, cellDataOld, err := v.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo old err: %s[%s]", err.Error(), k)
			return
		}
		n, ok := req.TxDidCellMap.Outputs[k]
		if !ok {
			resp.Err = fmt.Errorf("TxDidCellMap diff err: %s[%s]", err.Error(), k)
			return
		}
		_, cellDataNew, err := n.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo new err: %s[%s]", err.Error(), k)
			return
		}
		account := cellDataOld.Account
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
		mode := address.Mainnet
		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		addrOld, err := address.ConvertScriptToAddress(mode, v.Lock)
		if err != nil {
			resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
			return
		}
		oldOutpoint := common.OutPointStruct2String(v.OutPoint)
		oldOutpointList = append(oldOutpointList, oldOutpoint)

		tmp := tables.TableDidCellInfo{
			BlockNumber:  req.BlockNumber,
			Outpoint:     common.OutPointStruct2String(n.OutPoint),
			AccountId:    accountId,
			Account:      account,
			Args:         common.Bytes2Hex(n.Lock.Args),
			LockCodeHash: n.Lock.CodeHash.Hex(),
			ExpiredAt:    cellDataNew.ExpireAt,
		}
		list = append(list, tmp)

		if !v.Lock.Equals(n.Lock) {
			listTx = append(listTx, tables.TableTxInfo{
				Outpoint:       common.OutPointStruct2String(n.OutPoint),
				Action:         common.DidCellActionEditOwner,
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				AccountId:      accountId,
				Account:        account,
				Args:           common.Bytes2Hex(v.Lock.Args),
				Address:        addrOld,
				LockCodeHash:   v.Lock.CodeHash.Hex(),
			})
		}
		if cellDataOld.ExpireAt != cellDataNew.ExpireAt {
			listTx = append(listTx, tables.TableTxInfo{
				Outpoint:       common.OutPointStruct2String(n.OutPoint),
				Action:         common.DidCellActionRenew,
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				AccountId:      accountId,
				Account:        account,
				Args:           common.Bytes2Hex(v.Lock.Args),
				Address:        addrOld,
				LockCodeHash:   v.Lock.CodeHash.Hex(),
			})
		}
		if bytes.Compare(cellDataOld.WitnessHash, cellDataNew.WitnessHash) != 0 {
			listTx = append(listTx, tables.TableTxInfo{
				Outpoint:       common.OutPointStruct2String(n.OutPoint),
				Action:         common.DidCellActionEditRecords,
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				AccountId:      accountId,
				Account:        account,
				Args:           common.Bytes2Hex(v.Lock.Args),
				Address:        addrOld,
				LockCodeHash:   v.Lock.CodeHash.Hex(),
			})
			accountIds = append(accountIds, accountId)
			if w, yes := txDidEntityWitness.Outputs[v.Index]; yes {
				for _, r := range w.DidCellWitnessDataV0.Records {
					records = append(records, tables.TableRecordsInfo{
						AccountId: accountId,
						Account:   account,
						Key:       r.Key,
						Type:      r.Type,
						Label:     r.Label,
						Value:     r.Value,
						Ttl:       strconv.FormatUint(uint64(r.TTL), 10),
					})
				}
			}
		}
	}

	if err := b.DbDao.DidCellUpdateList(oldOutpointList, list, accountIds, records, listTx); err != nil {
		resp.Err = fmt.Errorf("DidCellUpdateList err: %s", err.Error())
		return
	}

	return
}
func (b *BlockParser) DidCellActionRecycle(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {

	log.Info("DidCellActionRecycle:", req.BlockNumber, req.TxHash, req.Action)
	var oldOutpointList []string
	var accountIds []string
	var listTx []tables.TableTxInfo
	for k, v := range req.TxDidCellMap.Inputs {
		oldOutpoint := common.OutPointStruct2String(v.OutPoint)
		oldOutpointList = append(oldOutpointList, oldOutpoint)

		_, cellData, err := v.GetDataInfo()
		if err != nil {
			resp.Err = fmt.Errorf("GetDataInfo err: %s[%s]", err.Error(), k)
			return
		}
		account := cellData.Account
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
		accountIds = append(accountIds, accountId)

		mode := address.Mainnet
		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		addr, err := address.ConvertScriptToAddress(mode, v.Lock)
		if err != nil {
			resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
			return
		}
		tmpTx := tables.TableTxInfo{
			Outpoint:       common.OutPoint2String(req.TxHash, v.OutPoint.Index),
			BlockNumber:    req.BlockNumber,
			BlockTimestamp: req.BlockTimestamp,
			AccountId:      accountId,
			Account:        account,
			Action:         req.Action,
			Args:           common.Bytes2Hex(v.Lock.Args),
			Address:        addr,
			LockCodeHash:   v.Lock.CodeHash.Hex(),
		}
		listTx = append(listTx, tmpTx)
	}

	if err := b.DbDao.DidCellRecycleList(oldOutpointList, accountIds, listTx); err != nil {
		resp.Err = fmt.Errorf("DidCellRecycleList err: %s", err.Error())
		return
	}
	return
}
