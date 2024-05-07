package handle

import (
	"did-cell-indexer/cache"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqTxSend struct {
	SignKey     string               `json:"sign_key"`
	SignAddress string               `json:"sign_address"`
	SignList    []txbuilder.SignData `json:"sign_list"`
}

type RespTxSend struct {
	Hash string `json:"hash"`
}

func (h *HttpHandle) TxSend(ctx *gin.Context) {
	var (
		funcName             = "TxSend"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqTxSend
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

	if err = h.doTxSend(&req, &apiResp); err != nil {
		log.Error("doTxSend err:", err.Error(), funcName, clientIp, remoteAddr)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doTxSend(req *ReqTxSend, apiResp *http_api.ApiResp) error {
	var resp RespTxSend

	var sic cache.SignInfoCache
	// get tx by cache
	if txStr, err := h.RC.GetSignTxCache(req.SignKey); err != nil {
		if err == redis.Nil {
			apiResp.ApiRespErr(http_api.ApiCodeTxExpired, "tx expired err")
		} else {
			apiResp.ApiRespErr(http_api.ApiCodeCacheError, "cache err")
		}
		return fmt.Errorf("GetSignTxCache err: %s", err.Error())
	} else if err = json.Unmarshal([]byte(txStr), &sic); err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "json.Unmarshal err")
		return fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}

	// sign
	txBuilder := txbuilder.NewDasTxBuilderFromBase(h.TxBuilderBase, sic.BuilderTx)
	if err := txBuilder.AddSignatureForTx(req.SignList); err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "add signature fail")
		return fmt.Errorf("AddSignatureForTx err: %s", err.Error())
	}

	// send tx
	hash, err := txBuilder.SendTransaction()
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "send tx err:"+err.Error())
		return fmt.Errorf("SendTransaction err: %s", err.Error())
	}
	resp.Hash = hash.Hex()
	// cache
	var outpoints []string
	for _, v := range txBuilder.Transaction.Inputs {
		outpoints = append(outpoints, common.OutPoint2String(v.PreviousOutput.TxHash.String(), v.PreviousOutput.Index))
	}
	h.DasCache.AddOutPoint(outpoints)
	// .bit balance ckb pay
	if sic.Action == common.DasActionTransfer {
		apiResp.ApiRespOK(resp)
		return nil
	}

	apiResp.ApiRespOK(resp)
	return nil
}
