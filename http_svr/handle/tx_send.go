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
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqTxSend struct {
	SignInfo
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
	if req.CKBTx != "" {
		log.Info("CKBTx:", req.CKBTx)
		userTx, err := rpc.TransactionFromString(req.CKBTx)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, fmt.Sprintf("rpc.TransactionFromString err: %s", err.Error()))
			return fmt.Errorf("rpc.TransactionFromString err: %s", err.Error())
		}
		cacheTxHash, err := sic.BuilderTx.Transaction.ComputeHash()
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, "ComputeHash err")
			return fmt.Errorf("ComputeHash err: %s", err.Error())
		}
		userTxHash, err := userTx.ComputeHash()
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, "ComputeHash err")
			return fmt.Errorf("ComputeHash err: %s", err.Error())
		}
		if userTxHash.String() != cacheTxHash.String() {
			apiResp.ApiRespErr(http_api.ApiCodeError500, "ckb tx invalid")
			return nil
		}
		sic.BuilderTx.Transaction = userTx
	}

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

	apiResp.ApiRespOK(resp)
	return nil
}
