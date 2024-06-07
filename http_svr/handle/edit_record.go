package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ReqEditRecord struct {
	CkbAddr  string `json:"ckb_addr" binding:"required" `
	Account  string `json:"account" binding:"required" binding:"required"`
	RawParam struct {
		Records []ReqRecord `json:"records"`
	} `json:"raw_param" binding:"required"`
}

type ReqRecord struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

type RespEditRecord struct {
	SignInfo
}

func (h *HttpHandle) RpcEditRecord(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqEditRecord
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	} else if len(req) == 0 {
		log.Error("len(req) is 0")
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doEditRecord(&req[0], apiResp); err != nil {
		log.Error("doEditRecord err:", err.Error())
	}
}

func (h *HttpHandle) EditRecord(ctx *gin.Context) {
	var (
		funcName             = "EditRecord"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqEditRecord
		apiResp              http_api.ApiResp
		err                  error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, ctx, remoteAddr)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req), ctx)

	if err = h.doEditRecord(&req, &apiResp); err != nil {
		log.Error("doEditRecord err:", err.Error(), funcName, clientIp, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doEditRecord(req *ReqEditRecord, apiResp *http_api.ApiResp) error {
	var resp RespEditRecord
	parseAddr, err := address.Parse(req.CkbAddr)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "ckb_addr error")
		log.Warnf("address.Parse err: %s", err.Error())
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	//args := common.Bytes2Hex(parseAddr.Script.Args)
	if req.Account == "" {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "account is invalid")
		return nil
	}
	if err := h.checkSystemUpgrade(apiResp); err != nil {
		return fmt.Errorf("checkSystemUpgrade err: %s", err.Error())
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	acc, err := h.DbDao.GetAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search account err")
		return fmt.Errorf("SearchAccount err: %s", err.Error())
	}
	if acc.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "account not exist")
		return nil
	} else if acc.IsExpired() {
		apiResp.ApiRespErr(http_api.ApiCodeAccountIsExpired, "account is expired")
		return nil
	} else if bytes.Compare(common.Hex2Bytes(acc.Args), parseAddr.Script.Args) != 0 {
		apiResp.ApiRespErr(http_api.ApiCodeNoAccountPermissions, "transfer account permission denied")
		return nil
	}
	outpoint := common.String2OutPointStruct(acc.Outpoint)
	var records []witness.Record
	builder, err := h.DasCore.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsRecordNamespace)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
		return fmt.Errorf("ConfigCellDataBuilderByTypeArgsList err: %s", err.Error())
	}
	log.Info("ConfigCellRecordKeys:", builder.ConfigCellRecordKeys)
	var mapRecordKey = make(map[string]struct{})
	for _, v := range builder.ConfigCellRecordKeys {
		mapRecordKey[v] = struct{}{}
	}

	for _, v := range req.RawParam.Records {
		record := fmt.Sprintf("%s.%s", v.Type, v.Key)
		if v.Type == "custom_key" { // (^[0-9a-z_]+$)
			if ok, _ := regexp.MatchString("^[0-9a-z_]+$", v.Key); !ok {
				apiResp.ApiRespErr(http_api.ApiCodeRecordInvalid, fmt.Sprintf("record [%s] is invalid", record))
				return nil
			}
		} else if v.Type == "address" {
			if ok, _ := regexp.MatchString("^(0|[1-9][0-9]*)$", v.Key); !ok {
				if _, ok2 := mapRecordKey[record]; !ok2 {
					apiResp.ApiRespErr(http_api.ApiCodeRecordInvalid, fmt.Sprintf("record [%s] is invalid", record))
					return nil
				}
			}
		} else if _, ok := mapRecordKey[record]; !ok {
			apiResp.ApiRespErr(http_api.ApiCodeRecordInvalid, fmt.Sprintf("record [%s] is invalid", record))
			return nil
		}
		ttl, err := strconv.ParseInt(v.TTL, 10, 64)
		if err != nil {
			ttl = 300
		}
		records = append(records, witness.Record{
			Key:   v.Key,
			Type:  v.Type,
			Label: v.Label,
			Value: v.Value,
			TTL:   uint32(ttl),
		})
	}

	recordsMolecule := witness.ConvertToCellRecords(records)
	recordsBys := recordsMolecule.AsSlice()
	log.Info("doEditRecord recordsBys:", len(recordsBys))
	if len(recordsBys) >= 5000 {
		apiResp.ApiRespErr(http_api.ApiCodeTooManyRecords, "too many records")
		return nil
	}

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
		Action:          common.DidCellActionEditRecords,
		DidCellOutPoint: outpoint,

		EditRecords: records,
	})
	if err != nil {
		log.Error("txbuilder.BuildDidCellTx err : ", err.Error())
		return fmt.Errorf("buildEditManagerTx err: %s", err.Error())
	}
	reqBuild := reqBuildTx{
		Action:  common.DidCellActionEditRecords,
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

func checkBuildTxErr(err error, apiResp *http_api.ApiResp) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "not live") {
		apiResp.ApiRespErr(http_api.ApiCodeOperationFrequent, "the operation is too frequent")
	} else {
		apiResp.ApiRespErr(http_api.ApiCodeError500, "Failed to build tx")
	}
}
