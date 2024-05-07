package handle

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

type ReqAccountRecords struct {
	Account string `json:"account"`
}

type RespAccountRecords struct {
	Records []RespAccountRecordsData `json:"records"`
}

type RespAccountRecordsData struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	Ttl   string `json:"ttl"`
}

func (h *HttpHandle) RpcAccountRecords(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqAccountRecords
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

	if err = h.doAccountRecords(&req[0], apiResp); err != nil {
		log.Error("doVersion err:", err.Error())
	}
}

func (h *HttpHandle) AccountRecords(ctx *gin.Context) {
	var (
		funcName             = "AccountRecords"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqAccountRecords
		apiResp              http_api.ApiResp
		err                  error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp, remoteAddr, ctx)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req), ctx)

	if err = h.doAccountRecords(&req, &apiResp); err != nil {
		log.Error("doAccountRecords err:", err.Error(), funcName, clientIp, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountRecords(req *ReqAccountRecords, apiResp *http_api.ApiResp) error {
	var resp RespAccountRecords
	resp.Records = make([]RespAccountRecordsData, 0)

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
		apiResp.ApiRespOK(resp)
		return nil
	}

	list, err := h.DbDao.SearchRecordsByAccount(accountId)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search records err")
		return fmt.Errorf("SearchRecordsByAccount err: %s", err.Error())
	}
	for _, v := range list {
		resp.Records = append(resp.Records, RespAccountRecordsData{
			Key:   v.Key,
			Type:  v.Type,
			Label: v.Label,
			Value: v.Value,
			Ttl:   v.Ttl,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}
