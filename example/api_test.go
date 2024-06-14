package example

import (
	"did-cell-indexer/http_svr/handle"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/parnurzeal/gorequest"
	"github.com/scorpiotzh/toolib"
	"testing"
)

var (
	TestUrl = "https://test-didcell-api.d.id/v1"
)

func doTxSend(t *testing.T, signInfo handle.SignInfo) {
	req := handle.ReqTxSend{SignInfo: signInfo}
	url := TestUrl + "/tx/send"
	var data handle.RespTxSend
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")
}

func TestRecycle(t *testing.T) {
	req := handle.ReqRecycle{
		ChainTypeAddress: core.ChainTypeAddress{
			Type: "blockchain",
			KeyInfo: core.KeyInfo{
				CoinType: common.CoinTypeCKB,
				Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
			},
		},
		Account: "20230616.bit",
	}
	url := TestUrl + "/recycle"
	var data handle.RespRecycle
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")

	//doTxSend(t, data.SignInfo)
}

func TestEditRecord(t *testing.T) {
	req := handle.ReqEditRecord{
		ChainTypeAddress: core.ChainTypeAddress{
			Type: "blockchain",
			KeyInfo: core.KeyInfo{
				CoinType: common.CoinTypeCKB,
				Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
			},
		},
		Account: "20230616.bit",
		Records: []handle.ReqRecord{{
			Key:   "309",
			Type:  "address",
			Label: "",
			Value: "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
			TTL:   "",
		}},
	}
	url := TestUrl + "/edit/record"
	var data handle.RespEditRecord
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")

	//doTxSend(t, data.SignInfo)
}

func TestEditOwner(t *testing.T) {
	req := handle.ReqEditOwner{
		ChainTypeAddress: core.ChainTypeAddress{
			Type: "blockchain",
			KeyInfo: core.KeyInfo{
				CoinType: common.CoinTypeCKB,
				Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
			},
		},
		Account:        "20230616.bit",
		ReceiveCkbAddr: "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgytmmrfg7aczevlxngqnr28npj2849erjyqqhe2guh",
	}
	url := TestUrl + "/edit/owner"
	var data handle.RespEditOwner
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")

	//doTxSend(t, data.SignInfo)
}

func TestRecordList(t *testing.T) {
	req := handle.ReqRecordList{
		Account: "20230616.bit",
	}
	url := TestUrl + "/record/list"
	var data handle.RespRecordList
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")
}

func TestAccountList(t *testing.T) {
	req := handle.ReqAccountList{
		ChainTypeAddress: core.ChainTypeAddress{
			Type: "blockchain",
			KeyInfo: core.KeyInfo{
				CoinType: common.CoinTypeCKB,
				Key:      "ckt1qyqrekdjpy72kvhp3e9uf6y5868w5hjg8qnsqt6a0m",
			},
		},
		Pagination: handle.Pagination{
			Page: 1,
			Size: 20,
		},
		Keyword: "",
		DidType: 0,
	}
	url := TestUrl + "/account/list"
	var data handle.RespAccountList
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
	fmt.Println("===========================")
}

func doReq(url string, req, data interface{}) error {
	var resp http_api.ApiResp
	resp.Data = &data

	_, _, errs := gorequest.New().Post(url).SendStruct(&req).EndStruct(&resp)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return fmt.Errorf("%d - %s", resp.ErrNo, resp.ErrMsg)
	}
	return nil
}
