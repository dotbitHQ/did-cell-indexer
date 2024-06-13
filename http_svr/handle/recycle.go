package handle

import (
	"did-cell-indexer/config"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqRecycle struct {
	core.ChainTypeAddress
	Account string `json:"account" binding:"required"`
}

type RespRecycle struct {
	SignInfo
}

func (h *HttpHandle) Recycle(ctx *gin.Context) {
	var (
		funcName             = "Recycle"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqRecycle
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

	if err = h.doRecycle(&req, &apiResp); err != nil {
		log.Error("doTransfer err:", err.Error(), funcName, clientIp, remoteAddr)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doRecycle(req *ReqRecycle, apiResp *http_api.ApiResp) error {
	var resp RespRecycle

	addrHex, err := req.FormatChainTypeAddress(config.Cfg.Server.Net, true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	} else if addrHex.DasAlgorithmId != common.DasAlgorithmIdAnyLock {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return nil
	}
	args := common.Bytes2Hex(addrHex.ParsedAddress.Script.Args)
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))

	didAccount, err := h.DbDao.GetAccountInfoForRecycle(accountId, args)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "Failed to get did account info")
		return fmt.Errorf("GetAccountInfoForRecycle err: %s", err.Error())
	} else if didAccount.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "account not exist")
		return nil
	}

	expiredAt := tables.GetDidCellRecycleExpiredAt()
	if didAccount.ExpiredAt > expiredAt {
		apiResp.ApiRespErr(http_api.ApiCodeNotYetDueForRecycle, "not yet due for recycle")
		return nil
	}

	txParams, err := txbuilder.BuildDidCellTx(txbuilder.DidCellTxParams{
		DasCore:         h.DasCore,
		DasCache:        h.DasCache,
		Action:          common.DidCellActionRecycle,
		DidCellOutPoint: didAccount.GetOutpoint(),
	})
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "Failed to build recycle tx")
		return fmt.Errorf("BuildDidCellTx err: %s", err.Error())
	}

	reqBuild := reqBuildTx{
		Action:  common.DidCellActionRecycle,
		Address: req.KeyInfo.Key,
		Account: req.Account,
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
