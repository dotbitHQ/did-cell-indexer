package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

type ReqTransfer struct {
	Account        string `json:"account" binding:"required"`
	CkbAddr        string `json:"ckb_addr" binding:"required" `
	ReceiveCkbAddr string `json:"receive_ckb_addr" binding:"required"`
}

type RespTransfer struct {
	SignInfo
}

func (h *HttpHandle) Transfer(ctx *gin.Context) {
	var (
		funcName             = "Transfer"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqTransfer
		apiResp              http_api.ApiResp
		err                  error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, remoteAddr)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, remoteAddr, toolib.JsonString(req))

	if err = h.doTransfer(&req, &apiResp); err != nil {
		log.Error("doTransfer err:", err.Error(), funcName, clientIp, remoteAddr)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTransfer(req *ReqTransfer, apiResp *http_api.ApiResp) error {
	var resp RespTransfer
	parseAddr, err := address.Parse(req.CkbAddr)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "ckb_addr error")
		log.Warnf("address.Parse err: %s", err.Error())
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	args := common.Bytes2Hex(parseAddr.Script.Args)

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	acc, err := h.DbDao.GetAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search account err")
		return fmt.Errorf("SearchAccount err: %s", err.Error())
	}
	if acc.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	} else if acc.ExpiredAt <= uint64(time.Now().Unix()) {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "account has expired")
		return nil
	} else if acc.Args != args {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "no permission")
		return nil
	}
	outpoint := common.String2OutPointStruct(acc.Outpoint)
	//todo api code and  完整 log
	receiveParseAddr, err := address.Parse(req.ReceiveCkbAddr)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "receive_ckb_addr error")
		log.Warnf("address.Parse err: %s", err.Error())
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	//receiveArgs := common.Bytes2Hex(receiveParseAddr.Script.Args)

	//fee cell
	//_, liveBalanceCell, err := h.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
	//	LockScript:   h.ServerScript,
	//	CapacityNeed: 5000,
	//	DasCache:     h.DasCache,
	//	SearchOrder:  indexer.SearchOrderDesc,
	//})
	//if err != nil {
	//	log.Warnf("GetBalanceCell err %s", err.Error())
	//	apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "GetBalanceCellWithLock error")
	//
	//	return fmt.Errorf("GetBalanceCell err %s", err.Error())
	//}

	txParams, err := txbuilder.BuildDidCellTx(txbuilder.DidCellTxParams{
		DasCore:         h.DasCore,
		DasCache:        h.DasCache,
		Action:          common.DidCellActionEditOwner,
		DidCellOutPoint: outpoint,
		EditOwnerLock:   receiveParseAddr.Script,
	})
	if err != nil {
		log.Error("txbuilder.BuildDidCellTx err : ", err.Error())
		return fmt.Errorf("buildEditManagerTx err: %s", err.Error())
	}
	reqBuild := reqBuildTx{
		Action:  common.DidCellActionEditOwner,
		Address: req.CkbAddr,
		Account: acc.Account,
	}
	if si, err := h.buildTx(&reqBuild, txParams); err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "build tx err")
		return fmt.Errorf("buildTx: %s", err.Error())
	} else {
		resp.SignInfo = *si
	}
	apiResp.ApiRespOK(resp)
	return nil
}
