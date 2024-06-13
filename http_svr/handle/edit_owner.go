package handle

import (
	"bytes"
	"did-cell-indexer/config"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqEditOwner struct {
	core.ChainTypeAddress
	Account        string `json:"account" binding:"required"`
	ReceiveCkbAddr string `json:"receive_ckb_addr" binding:"required"`
}

type RespEditOwner struct {
	SignInfo
}

func (h *HttpHandle) EditOwner(ctx *gin.Context) {
	var (
		funcName             = "EditOwner"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqEditOwner
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

	if err = h.doEditOwner(&req, &apiResp); err != nil {
		log.Error("doEditOwner err:", err.Error(), funcName, clientIp, remoteAddr)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doEditOwner(req *ReqEditOwner, apiResp *http_api.ApiResp) error {
	var resp RespEditOwner

	addrHexFrom, err := req.FormatChainTypeAddress(config.Cfg.Server.Net, true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address is invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	} else if addrHexFrom.DasAlgorithmId != common.DasAlgorithmIdAnyLock {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return nil
	}

	toCTA := core.ChainTypeAddress{
		Type: "blockchain",
		KeyInfo: core.KeyInfo{
			CoinType: common.CoinTypeCKB,
			Key:      req.ReceiveCkbAddr,
		},
	}
	addrHexTo, err := toCTA.FormatChainTypeAddress(config.Cfg.Server.Net, true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "receiver address is invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	} else if addrHexTo.DasAlgorithmId != common.DasAlgorithmIdAnyLock {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "receiver address is invalid")
		return nil
	}
	editOwnerLock := addrHexTo.ParsedAddress.Script

	if strings.EqualFold(req.KeyInfo.Key, req.ReceiveCkbAddr) {
		apiResp.ApiRespErr(http_api.ApiCodeSameLock, "same owner address")
		return nil
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	if err := h.checkSystemUpgrade(apiResp); err != nil {
		return fmt.Errorf("checkSystemUpgrade err: %s", err.Error())
	}

	acc, err := h.DbDao.GetAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "Failed to get account info")
		return fmt.Errorf("GetAccountInfoByAccountId err: %s", err.Error())
	} else if acc.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "account not exist")
		return nil
	} else if acc.IsExpired() {
		apiResp.ApiRespErr(http_api.ApiCodeAccountIsExpired, "account expired")
		return nil
	} else if bytes.Compare(common.Hex2Bytes(acc.Args), addrHexFrom.ParsedAddress.Script.Args) != 0 {
		apiResp.ApiRespErr(http_api.ApiCodeNoAccountPermissions, "transfer account permission denied")
		return nil
	}

	txParams, err := txbuilder.BuildDidCellTx(txbuilder.DidCellTxParams{
		DasCore:         h.DasCore,
		DasCache:        h.DasCache,
		Action:          common.DidCellActionEditOwner,
		DidCellOutPoint: acc.GetOutpoint(),
		EditOwnerLock:   editOwnerLock,
	})
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "Failed to build tx")
		return fmt.Errorf("BuildDidCellTx err: %s", err.Error())
	}

	reqBuild := reqBuildTx{
		Action:  common.DidCellActionEditOwner,
		Address: req.KeyInfo.Key,
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
