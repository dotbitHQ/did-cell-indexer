package handle

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

type ReqRecordList struct {
	Account string `json:"account"`
}

type RespRecordList struct {
	Records []RespRecordListData `json:"records"`
}

type RespRecordListData struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	Ttl   string `json:"ttl"`
}

func (h *HttpHandle) RecordList(ctx *gin.Context) {
	var (
		funcName             = "RecordList"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqRecordList
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

	if err = h.doRecordList(&req, &apiResp); err != nil {
		log.Error("doRecordList err:", err.Error(), funcName, clientIp, ctx)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doRecordList(req *ReqRecordList, apiResp *http_api.ApiResp) error {
	var resp RespRecordList
	resp.Records = make([]RespRecordListData, 0)

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
		resp.Records = append(resp.Records, RespRecordListData{
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
