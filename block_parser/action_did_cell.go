package block_parser

//
//func (b *BlockParser) ActionAccountUpgrade(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
//	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameAccountCellType); err != nil {
//		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
//		return
//	} else if !isCV {
//		log.Warn("not current version account cross chain tx")
//		return
//	}
//	log.Info("ActionAccountCrossChain:", req.BlockNumber, req.TxHash, req.Action)
//
//	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeOutputs)
//	if err != nil {
//		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
//		return
//	}
//	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[didEntity.Target.Index].Lock.Args)
//
//	account, expiredAt, err := witness.GetAccountAndExpireFromDidCellData(req.Tx.OutputsData[didEntity.Target.Index])
//	if err != nil {
//		resp.Err = fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//		return
//	}
//	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
//
//	didCellInfo := tables.TableDidCellInfo{
//		BlockNumber:  req.BlockNumber,
//		Outpoint:     common.OutPoint2String(req.TxHash, 0),
//		AccountId:    accountId,
//		Account:      account,
//		Args:         didCellArgs,
//		LockCodeHash: req.Tx.Outputs[didEntity.Target.Index].Lock.CodeHash.Hex(),
//		ExpiredAt:    expiredAt,
//	}
//
//	if err = b.DbDao.AccountUpgrade(didCellInfo); err != nil {
//		resp.Err = fmt.Errorf("AccountCrossChain err: %s ", err.Error())
//		return
//	}
//	return
//}

// edit record
//func (b *BlockParser) ActionEditDidCellRecords(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
//	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
//		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
//		return
//	} else if !isCV {
//		log.Warn("not current version account cross chain tx")
//		return
//	}
//	log.Info("ActionEditDidCellRecords:", req.BlockNumber, req.TxHash, req.Action)
//
//	txDidEntity, err := witness.TxToDidEntity(req.Tx)
//	if err != nil {
//		resp.Err = fmt.Errorf("TxToDidEntity err: %s", err.Error())
//		return
//	}
//
//	account, _, err := witness.GetAccountAndExpireFromDidCellData(req.Tx.OutputsData[txDidEntity.Outputs[0].Target.Index])
//	if err != nil {
//		resp.Err = fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//		return
//	}
//	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
//	log.Info("ActionEditDidRecords:", account)
//
//	var recordsInfos []tables.TableRecordsInfo
//	for _, v := range txDidEntity.Outputs[0].DidCellWitnessDataV0.Records {
//		recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
//			AccountId:       accountId,
//			ParentAccountId: "",
//			Account:         account,
//			Key:             v.Key,
//			Type:            v.Type,
//			Label:           v.Label,
//			Value:           v.Value,
//			Ttl:             strconv.FormatUint(uint64(v.TTL), 10),
//		})
//	}
//
//	oldOutpoint := common.OutPointStruct2String(req.Tx.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
//
//	var didCellInfo tables.TableDidCellInfo
//	didCellInfo.AccountId = accountId
//	didCellInfo.BlockNumber = req.BlockNumber
//	didCellInfo.Outpoint = common.OutPoint2String(req.Tx.Hash.Hex(), uint(txDidEntity.Outputs[0].Target.Index))
//
//	mode := address.Mainnet
//	if config.Cfg.Server.Net != common.DasNetTypeMainNet {
//		mode = address.Testnet
//	}
//	anyLockAddr, err := address.ConvertScriptToAddress(mode, req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock)
//	if err != nil {
//		resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
//		return
//	}
//
//	txInfo := tables.TableTxInfo{
//		Outpoint:       didCellInfo.Outpoint,
//		BlockNumber:    req.BlockNumber,
//		BlockTimestamp: req.BlockTimestamp,
//		AccountId:      accountId,
//		Account:        account,
//		Action:         common.DidCellActionEditRecords,
//		Args:           common.Bytes2Hex(req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.Args),
//		Address:        anyLockAddr,
//		LockCodeHash:   req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.CodeHash.Hex(),
//	}
//
//	if err := b.DbDao.CreateDidCellRecordsInfos(oldOutpoint, didCellInfo, recordsInfos, txInfo); err != nil {
//		resp.Err = fmt.Errorf("CreateDidCellRecordsInfos err: %s", err.Error())
//		return
//	}
//
//	return
//}
//
//func (b *BlockParser) ActionEditDidCellOwner(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
//	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
//		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
//		return
//	} else if !isCV {
//		log.Warn("not current version account cross chain tx")
//		return
//	}
//	log.Info("ActionEditDidCellOwner:", req.BlockNumber, req.TxHash, req.Action)
//
//	txDidEntity, err := witness.TxToDidEntity(req.Tx)
//	if err != nil {
//		resp.Err = fmt.Errorf("TxToDidEntity err: %s", err.Error())
//		return
//	}
//
//	account, expiredAt, err := witness.GetAccountAndExpireFromDidCellData(req.Tx.OutputsData[txDidEntity.Outputs[0].Target.Index])
//	if err != nil {
//		resp.Err = fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//		return
//	}
//
//	didCellArgs := common.Bytes2Hex(req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.Args)
//	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
//	didCellInfo := tables.TableDidCellInfo{
//		BlockNumber:  req.BlockNumber,
//		Outpoint:     common.OutPoint2String(req.TxHash, uint(txDidEntity.Outputs[0].Target.Index)),
//		AccountId:    accountId,
//		Account:      account,
//		Args:         didCellArgs,
//		LockCodeHash: req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.CodeHash.Hex(),
//		ExpiredAt:    expiredAt,
//	}
//
//	var recordsInfos []tables.TableRecordsInfo
//	for _, v := range txDidEntity.Outputs[0].DidCellWitnessDataV0.Records {
//		recordsInfos = append(recordsInfos, tables.TableRecordsInfo{
//			AccountId:       accountId,
//			ParentAccountId: "",
//			Account:         account,
//			Key:             v.Key,
//			Type:            v.Type,
//			Label:           v.Label,
//			Value:           v.Value,
//			Ttl:             strconv.FormatUint(uint64(v.TTL), 10),
//		})
//	}
//
//	txInfo := tables.TableTxInfo{
//		Outpoint:       didCellInfo.Outpoint,
//		BlockNumber:    req.BlockNumber,
//		BlockTimestamp: req.BlockTimestamp,
//		AccountId:      didCellInfo.AccountId,
//		Account:        didCellInfo.Account,
//		Action:         common.DidCellActionEditOwner,
//		Args:           "",
//		LockCodeHash:   "",
//		Address:        "",
//	}
//	var oldOutpoint string
//	if len(txDidEntity.Inputs) > 0 {
//		oldOutpoint = common.OutPointStruct2String(req.Tx.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
//		preInput := req.Tx.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput
//		preTx, err := b.DasCore.Client().GetTransaction(context.Background(), preInput.TxHash)
//		if err != nil {
//			resp.Err = fmt.Errorf("GetTransaction err: %s", err.Error())
//			return
//		}
//		txInfo.LockCodeHash = preTx.Transaction.Outputs[preInput.Index].Lock.CodeHash.Hex()
//		txInfo.Args = common.Bytes2Hex(preTx.Transaction.Outputs[preInput.Index].Lock.Args)
//
//		mode := address.Mainnet
//		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
//			mode = address.Testnet
//		}
//		anyLockAddr, err := address.ConvertScriptToAddress(mode, preTx.Transaction.Outputs[preInput.Index].Lock)
//		if err != nil {
//			resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
//			return
//		}
//		txInfo.Address = anyLockAddr
//	}
//
//	if err := b.DbDao.EditDidCellOwner(oldOutpoint, didCellInfo, recordsInfos, txInfo); err != nil {
//		resp.Err = fmt.Errorf("EditDidCellOwner err: %s", err.Error())
//		return
//	}
//	return
//}

//func (b *BlockParser) ActionDidCellRenew(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
//	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameDidCellType); err != nil {
//		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
//		return
//	} else if !isCV {
//		log.Warn("not current version renew account tx")
//		return
//	}
//	log.Info("ActionDidCellRenew:", req.BlockNumber, req.TxHash)
//
//	//renew did cell
//	txDidEntity, err := witness.TxToDidEntity(req.Tx)
//	if err != nil {
//		resp.Err = fmt.Errorf("TxToDidEntity err: %s", err.Error())
//		return
//	}
//
//	account, expiredAt, err := witness.GetAccountAndExpireFromDidCellData(req.Tx.OutputsData[txDidEntity.Outputs[0].Target.Index])
//	if err != nil {
//		resp.Err = fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//		return
//	}
//	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
//
//	var oldOutpoint string
//	if len(txDidEntity.Inputs) > 0 {
//		oldOutpoint = common.OutPointStruct2String(req.Tx.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
//	}
//
//	var didCellInfo tables.TableDidCellInfo
//	didCellInfo.Outpoint = common.OutPoint2String(req.TxHash, uint(txDidEntity.Outputs[0].Target.Index))
//	didCellInfo.ExpiredAt = expiredAt
//	didCellInfo.BlockNumber = req.BlockNumber
//
//	mode := address.Mainnet
//	if config.Cfg.Server.Net != common.DasNetTypeMainNet {
//		mode = address.Testnet
//	}
//	anyLockAddr, err := address.ConvertScriptToAddress(mode, req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock)
//	if err != nil {
//		resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
//		return
//	}
//	txInfo := tables.TableTxInfo{
//		Outpoint:       didCellInfo.Outpoint,
//		BlockNumber:    req.BlockNumber,
//		BlockTimestamp: req.BlockTimestamp,
//		AccountId:      accountId,
//		Account:        account,
//		Action:         common.DidCellActionRenew,
//		Args:           common.Bytes2Hex(req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.Args),
//		Address:        anyLockAddr,
//		LockCodeHash:   req.Tx.Outputs[txDidEntity.Outputs[0].Target.Index].Lock.CodeHash.Hex(),
//	}
//
//	if err := b.DbDao.DidCellRenew(oldOutpoint, didCellInfo, txInfo); err != nil {
//		resp.Err = fmt.Errorf("DidCellRenew err: %s", err.Error())
//		return
//	}
//	return
//}

//func (b *BlockParser) ActionDidCellRecycle(req FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
//	didEntity, err := witness.TxToOneDidEntity(req.Tx, witness.SourceTypeInputs)
//	if err != nil {
//		resp.Err = fmt.Errorf("TxToOneDidEntity err: %s", err.Error())
//		return
//	}
//
//	preTx, err := b.DasCore.Client().GetTransaction(b.Ctx, req.Tx.Inputs[didEntity.Target.Index].PreviousOutput.TxHash)
//	if err != nil {
//		resp.Err = fmt.Errorf("GetTransaction err: %s", err.Error())
//		return
//	}
//
//	if isCV, err := isCurrentVersionTx(preTx.Transaction, common.DasContractNameDidCellType); err != nil {
//		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err.Error())
//		return
//	} else if !isCV {
//		log.Warn("not current version didcell recycle")
//		return
//	}
//	log.Info("ActionDidCellRecycle:", req.BlockNumber, req.TxHash, req.Action)
//
//	account, _, err := witness.GetAccountAndExpireFromDidCellData(preTx.Transaction.OutputsData[req.Tx.Inputs[didEntity.Target.Index].PreviousOutput.Index])
//	if err != nil {
//		resp.Err = fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//		return
//	}
//
//	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
//	oldOutpoint := common.OutPointStruct2String(req.Tx.Inputs[didEntity.Target.Index].PreviousOutput)
//
//	mode := address.Mainnet
//	if config.Cfg.Server.Net != common.DasNetTypeMainNet {
//		mode = address.Testnet
//	}
//	anyLockAddr, err := address.ConvertScriptToAddress(mode, preTx.Transaction.Outputs[req.Tx.Inputs[didEntity.Target.Index].PreviousOutput.Index].Lock)
//	if err != nil {
//		resp.Err = fmt.Errorf("address.ConvertScriptToAddress err: %s", err.Error())
//		return
//	}
//	txInfo := tables.TableTxInfo{
//		Outpoint:       common.OutPoint2String(req.TxHash, 0),
//		BlockNumber:    req.BlockNumber,
//		BlockTimestamp: req.BlockTimestamp,
//		AccountId:      accountId,
//		Account:        account,
//		Action:         common.DidCellActionEditOwner,
//		Args:           common.Bytes2Hex(preTx.Transaction.Outputs[req.Tx.Inputs[didEntity.Target.Index].PreviousOutput.Index].Lock.Args),
//		Address:        anyLockAddr,
//		LockCodeHash:   preTx.Transaction.Outputs[req.Tx.Inputs[didEntity.Target.Index].PreviousOutput.Index].Lock.CodeHash.Hex(),
//	}
//
//	if err := b.DbDao.DidCellRecycle(oldOutpoint, accountId, txInfo); err != nil {
//		resp.Err = fmt.Errorf("DidCellRecycle err: %s", err.Error())
//		return
//	}
//	return
//
//}
