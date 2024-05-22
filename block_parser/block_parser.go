package block_parser

import (
	"context"
	"did-cell-indexer/config"
	"did-cell-indexer/dao"
	"did-cell-indexer/notify"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sync"
	"sync/atomic"
	"time"
)

var log = logger.NewLogger("block_parser", logger.LevelDebug)

type BlockParser struct {
	parserType           tables.ParserType
	mapTransactionHandle map[common.DasAction]FuncTransactionHandle
	DasCore              *core.DasCore
	CurrentBlockNumber   uint64
	DbDao                *dao.DbDao
	ConcurrencyNum       uint64
	ConfirmNum           uint64
	Ctx                  context.Context
	Wg                   *sync.WaitGroup
}

func (b *BlockParser) Run() error {
	b.parserType = tables.ParserTypeCKB
	b.registerTransactionHandle()
	currentBlockNumber, err := b.DasCore.Client().GetTipBlockNumber(b.Ctx)
	if err != nil {
		return fmt.Errorf("GetTipBlockNumber err: %s", err.Error())
	}

	if err := b.initCurrentBlockNumber(currentBlockNumber); err != nil {
		return fmt.Errorf("initCurrentBlockNumber err: %s", err.Error())
	}
	atomic.AddUint64(&b.CurrentBlockNumber, 1)
	b.Wg.Add(1)
	go func() {
		defer http_api.RecoverPanic()
		for {
			select {
			default:
				latestBlockNumber, err := b.DasCore.Client().GetTipBlockNumber(b.Ctx)
				if err != nil {
					log.Error("GetTipBlockNumber err:", err.Error())
				} else {
					if b.ConcurrencyNum > 1 && b.CurrentBlockNumber < (latestBlockNumber-b.ConfirmNum-b.ConcurrencyNum) {
						nowTime := time.Now()
						if err = b.parserConcurrencyMode(); err != nil {
							log.Error("parserConcurrencyMode err:", err.Error(), b.CurrentBlockNumber)
						}
						log.Debug("parserConcurrencyMode time:", time.Since(nowTime).Seconds())
					} else if b.CurrentBlockNumber < (latestBlockNumber - b.ConfirmNum) { // check rollback
						nowTime := time.Now()
						if err = b.parserSubMode(); err != nil {
							log.Error("parserSubMode err:", err.Error(), b.CurrentBlockNumber)
						}
						log.Debug("parserSubMode time:", time.Since(nowTime).Seconds())
					} else {
						log.Debug("RunParser:", b.CurrentBlockNumber, latestBlockNumber)
						time.Sleep(time.Second * 10)
					}
					time.Sleep(time.Millisecond * 300)
				}
			case <-b.Ctx.Done():
				b.Wg.Done()
				return
			}
		}
	}()
	return nil
}

func (b *BlockParser) initCurrentBlockNumber(currentBlockNumber uint64) error {
	if block, err := b.DbDao.FindBlockInfo(b.parserType); err != nil {
		return err
	} else if block.Id > 0 {
		b.CurrentBlockNumber = block.BlockNumber
	} else if b.CurrentBlockNumber == 0 && currentBlockNumber > 0 {
		b.CurrentBlockNumber = currentBlockNumber
	}
	return nil
}

func (b *BlockParser) parserSubMode() error {
	log.Debug("parserSubMode:", b.CurrentBlockNumber)
	block, err := b.DasCore.Client().GetBlockByNumber(b.Ctx, b.CurrentBlockNumber)
	if err != nil {
		return fmt.Errorf("GetBlockByNumber err: %s", err.Error())
	} else {
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Debug("parserSubMode:", b.CurrentBlockNumber, blockHash, parentHash)
		if fork, err := b.checkFork(parentHash); err != nil {
			return fmt.Errorf("checkFork err: %s", err.Error())
		} else if fork {
			log.Debug("CheckFork is true:", b.CurrentBlockNumber, blockHash, parentHash)
			atomic.AddUint64(&b.CurrentBlockNumber, ^uint64(0))
		} else if err = b.parsingBlockData(block); err != nil {
			return fmt.Errorf("parsingBlockData err: %s", err.Error())
		} else {
			if err = b.DbDao.CreateBlockInfo(b.parserType, b.CurrentBlockNumber, blockHash, parentHash); err != nil {
				return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
			} else {
				atomic.AddUint64(&b.CurrentBlockNumber, 1)
			}
			if err = b.DbDao.DeleteBlockInfo(b.parserType, b.CurrentBlockNumber-20); err != nil {
				return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
			}
		}
	}
	return nil
}

func (b *BlockParser) checkFork(parentHash string) (bool, error) {
	block, err := b.DbDao.FindBlockInfoByBlockNumber(b.parserType, b.CurrentBlockNumber-1)
	if err != nil {
		return false, err
	}
	if block.Id == 0 {
		return false, nil
	} else if block.BlockHash != parentHash {
		log.Warn("CheckFork:", b.CurrentBlockNumber, parentHash, block.BlockHash)
		return true, nil
	}
	return false, nil
}

func (b *BlockParser) parserConcurrencyMode() error {
	log.Debug("parserConcurrencyMode:", b.CurrentBlockNumber, b.ConcurrencyNum)
	for i := uint64(0); i < b.ConcurrencyNum; i++ {
		block, err := b.DasCore.Client().GetBlockByNumber(b.Ctx, b.CurrentBlockNumber)
		if err != nil {
			return fmt.Errorf("GetBlockByNumber err: %s [%d]", err.Error(), b.CurrentBlockNumber)
		}
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Debug("parserConcurrencyMode:", b.CurrentBlockNumber, blockHash, parentHash)

		if err = b.parsingBlockData(block); err != nil {
			return fmt.Errorf("parsingBlockData err: %s", err.Error())
		} else {
			if err = b.DbDao.CreateBlockInfo(b.parserType, b.CurrentBlockNumber, blockHash, parentHash); err != nil {
				return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
			} else {
				atomic.AddUint64(&b.CurrentBlockNumber, 1)
			}
		}
	}
	if err := b.DbDao.DeleteBlockInfo(b.parserType, b.CurrentBlockNumber-20); err != nil {
		return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
	}
	return nil
}

func (b *BlockParser) parsingBlockData(block *types.Block) error {
	for _, tx := range block.Transactions {
		txHash := tx.Hash.Hex()
		blockNumber := block.Header.Number
		blockTimestamp := block.Header.Timestamp
		if didCellAction, err := b.DasCore.TxToDidCellAction(tx); err != nil {
			if err != core.ErrorNotExistDidCell {
				log.Error("TxToDidCellAction err :", didCellAction, blockNumber, txHash, err.Error())
			}
		} else {
			if handle, ok := b.mapTransactionHandle[didCellAction]; ok {
				// transaction parse by action
				resp := handle(FuncTransactionHandleReq{
					DbDao:          b.DbDao,
					Tx:             tx,
					TxHash:         txHash,
					BlockNumber:    blockNumber,
					BlockTimestamp: int64(blockTimestamp),
					Action:         didCellAction,
				})
				if resp.Err != nil {
					notify.SendLarkErrNotify(config.Cfg.Notify.LarkErrKey, "parsingBlockData", err.Error())
					log.Error("action handle resp:", didCellAction, blockNumber, txHash, resp.Err.Error())
					return resp.Err
				}
			}
		}

	}
	return nil
}
