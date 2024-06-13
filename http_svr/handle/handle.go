package handle

import (
	"context"
	"did-cell-indexer/cache"
	"did-cell-indexer/config"
	"did-cell-indexer/dao"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
)

var (
	log = logger.NewLoggerDefault("http_handle", logger.LevelDebug, nil)
)

type HttpHandle struct {
	Ctx           context.Context
	DbDao         *dao.DbDao
	RC            *cache.RedisCache
	DasCore       *core.DasCore
	DasCache      *dascache.DasCache
	ServerScript  *types.Script
	TxBuilderBase *txbuilder.DasTxBuilderBase
}

func GetClientIp(ctx *gin.Context) (string, string) {
	clientIP := fmt.Sprintf("%v", ctx.Request.Header.Get("X-Real-IP"))
	return clientIP, ctx.Request.RemoteAddr
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p Pagination) GetLimit() int {
	if p.Size < 1 || p.Size > 100 {
		return 100
	}
	return p.Size
}

func (p Pagination) GetOffset() int {
	page := p.Page
	if p.Page < 1 {
		page = 1
	}
	size := p.GetLimit()
	return (page - 1) * size
}

func (h *HttpHandle) checkSystemUpgrade(apiResp *http_api.ApiResp) error {
	if config.Cfg.Server.IsUpdate {
		apiResp.ApiRespErr(http_api.ApiCodeSystemUpgrade, http_api.TextSystemUpgrade)
		return nil
	}
	return nil
}

type reqBuildTx struct {
	Action  common.DasAction
	Address string `json:"address"`
	Account string `json:"account"`
}
type SignInfo struct {
	SignKey  string               `json:"sign_key"`  // sign tx key
	SignList []txbuilder.SignData `json:"sign_list"` // sign list
	CKBTx    string               `json:"ckb_tx"`
}

func (h *HttpHandle) buildTx(req *reqBuildTx, txParams *txbuilder.BuildTransactionParams) (*SignInfo, error) {

	txBuilder := txbuilder.NewDasTxBuilderFromBase(h.TxBuilderBase, nil)
	if err := txBuilder.BuildTransaction(txParams); err != nil {
		return nil, fmt.Errorf("txBuilder.BuildTransaction err: %s", err.Error())
	}
	sizeInBlock, _ := txBuilder.Transaction.SizeInBlock()
	txFeeRate := config.Cfg.Server.TxTeeRate
	if txFeeRate == 0 {
		txFeeRate = 1
	}
	txFee := txFeeRate*sizeInBlock + 1000
	log.Info("buildTx tx fee:", req.Action, txFee, sizeInBlock, txFee)
	var skipGroups []int

	switch req.Action {
	case common.DidCellActionRecycle,
		common.DidCellActionEditRecords,
		common.DidCellActionEditOwner:
		changeCapacity := txBuilder.Transaction.Outputs[0].Capacity - txFee
		txBuilder.Transaction.Outputs[0].Capacity = changeCapacity
		log.Info("buildTx user:", req.Action, sizeInBlock, changeCapacity)
	}

	signList, err := txBuilder.GenerateDigestListFromTx(skipGroups)
	if err != nil {
		return nil, fmt.Errorf("txBuilder.GenerateDigestListFromTx err: %s", err.Error())
	}

	txStr := txBuilder.TxString()
	log.Info("buildTx:", txStr)

	var sic cache.SignInfoCache
	sic.Action = req.Action
	sic.CkbAddr = req.Address
	sic.BuilderTx = txBuilder.DasTxBuilderTransaction

	signKey := sic.SignKey()
	cacheStr := toolib.JsonString(&sic)
	if err = h.RC.SetSignTxCache(signKey, cacheStr); err != nil {
		return nil, fmt.Errorf("SetSignTxCache err: %s", err.Error())
	}

	var si SignInfo
	si.SignKey = signKey
	si.SignList = signList
	si.CKBTx = txStr

	return &si, nil
}
