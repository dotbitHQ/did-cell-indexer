package handle

import (
	"context"
	"did-cell-indexer/config"
	"did-cell-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

type ReqAccountList struct {
	core.ChainTypeAddress
	Pagination
	Keyword string               `json:"keyword"`
	DidType tables.DidCellStatus `json:"did_type"`
}

type RespAccountList struct {
	Total int64     `json:"total"`
	List  []DidData `json:"did_list"`
}

type DidData struct {
	Outpoint      string               `json:"outpoint"`
	AccountId     string               `json:"account_id"`
	Account       string               `json:"account"`
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
	resp.List = make([]DidData, 0)

	addrHex, err := req.FormatChainTypeAddress(config.Cfg.Server.Net, true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	} else if addrHex.DasAlgorithmId != common.DasAlgorithmIdAnyLock {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return nil
	}
	args := common.Bytes2Hex(addrHex.ParsedAddress.Script.Args)

	res, err := h.DbDao.QueryDidCell(args, req.Keyword, req.GetLimit(), req.GetOffset(), req.DidType)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search did cell list err")
		return fmt.Errorf("QueryDidCell err: %s", err.Error())
	}

	recycleAt := tables.GetDidCellRecycleExpiredAt()
	nowAt := uint64(time.Now().Unix())
	for _, v := range res {
		temp := DidData{
			Outpoint:  v.Outpoint,
			Account:   v.Account,
			AccountId: v.AccountId,
			ExpiredAt: v.ExpiredAt,
		}

		if v.ExpiredAt >= nowAt {
			temp.DidCellStatus = tables.DidCellStatusNormal
		} else if v.ExpiredAt <= nowAt && v.ExpiredAt >= recycleAt {
			temp.DidCellStatus = tables.DidCellStatusExpired
		} else if v.ExpiredAt < recycleAt {
			temp.DidCellStatus = tables.DidCellStatusRecycle
		}
		resp.List = append(resp.List, temp)
	}

	resp.Total, err = h.DbDao.QueryDidCellTotal(args, req.Keyword, req.DidType)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search did cell count err")
		return fmt.Errorf("QueryDidCellTotal err: %s", err.Error())
	}

	apiResp.ApiRespOK(resp)
	return nil
}
