package handle

import (
	"context"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

type ReqAccountList struct {
	CkbAddr string               `json:"ckb_addr" binding:"required"`
	DidType tables.DidCellStatus `json:"did_type"`
}

type RespAccountList struct {
	List []DidData `json:"did_list"`
}

type DidData struct {
	Outpoint      string               `json:"outpoint"`
	AccountId     string               `json:"account_id"`
	Account       string               `json:"account"`
	Args          string               `json:"args"`
	ExpiredAt     uint64               `json:"expired_at"`
	DidCellStatus tables.DidCellStatus `json:"did_cell_status"`
}

func (h *HttpHandle) AccountList(ctx *gin.Context) {
	var (
		funcName             = "AccountList"
		clientIp, remoteAddr = GetClientIp(ctx)
		req                  ReqAccountList
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

	if err = h.doAccountList(ctx, &req, &apiResp); err != nil {
		log.Error(ctx, "doAccountList err:", err.Error(), funcName, clientIp, remoteAddr)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountList(ctx context.Context, req *ReqAccountList, apiResp *http_api.ApiResp) error {
	var resp RespAccountList
	data := make([]DidData, 0)
	parseAddr, err := address.Parse(req.CkbAddr)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "ckb address error")
		log.Warnf("address.Parse err: %s", err.Error())
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	args := common.Bytes2Hex(parseAddr.Script.Args)
	res, err := h.DbDao.QueryDidCell(args, req.DidType)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search account list err")
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	for _, v := range res {
		temp := DidData{
			Outpoint:  v.Outpoint,
			Account:   v.Account,
			AccountId: v.AccountId,
			Args:      v.Args,
			ExpiredAt: v.ExpiredAt,
		}
		if v.ExpiredAt > uint64(time.Now().Unix()) {
			temp.DidCellStatus = tables.DidCellStatusNormal
		} else {
			temp.DidCellStatus = tables.DidCellStatusExpired
		}
		data = append(data, temp)
	}
	resp.List = data
	apiResp.ApiRespOK(resp)
	return nil
}
