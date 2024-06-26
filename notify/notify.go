package notify

import (
	"did-cell-indexer/prometheus"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/parnurzeal/gorequest"
	"time"
)

var log = logger.NewLoggerDefault("notify", logger.LevelDebug, nil)

const (
	LarkNotifyUrl = "https://open.larksuite.com/open-apis/bot/v2/hook/%s"
)

type MsgContent struct {
	Tag      string `json:"tag"`
	UnEscape bool   `json:"un_escape"`
	Text     string `json:"text"`
	UserName string `json:"user_name"`
	UserId   string `json:"user_id"`
}
type MsgData struct {
	Email   string `json:"email"`
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string         `json:"title"`
				Content [][]MsgContent `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

func SendLarkErrNotify(key, title, text string) {
	prometheus.Tools.Metrics.ErrNotify().WithLabelValues(title, text).Inc()
	//SendLarkTextNotify(key, title, text)
}

func SendLarkTextNotify(key, title, text string) {
	if key == "" || text == "" {
		return
	}
	var data MsgData
	data.Email = ""
	data.MsgType = "post"
	data.Content.Post.ZhCn.Title = fmt.Sprintf("DP-SVR: %s", title)
	data.Content.Post.ZhCn.Content = [][]MsgContent{
		{
			MsgContent{
				Tag:      "text",
				UnEscape: false,
				Text:     text,
			},
		},
	}
	url := fmt.Sprintf(LarkNotifyUrl, key)
	_, body, errs := gorequest.New().Post(url).Timeout(time.Second*10).Retry(5, 5*time.Second).SendStruct(&data).End()
	if len(errs) > 0 {
		log.Error("SendLarkTextNotify req err:", errs)
	} else {
		log.Info("SendLarkTextNotify req:", body)
	}
}

func SendLarkTextNotifyAtAll(key, title, text string) {
	if key == "" || text == "" {
		return
	}
	var data MsgData
	data.Email = ""
	data.MsgType = "post"
	data.Content.Post.ZhCn.Title = fmt.Sprintf("DP-SVR: %s", title)
	data.Content.Post.ZhCn.Content = [][]MsgContent{
		{
			MsgContent{
				Tag:      "text",
				UnEscape: false,
				Text:     text,
			},
		},
		{
			MsgContent{
				Tag:      "at",
				UserId:   "all",
				UserName: "所有人",
			},
		},
	}
	url := fmt.Sprintf(LarkNotifyUrl, key)
	_, body, errs := gorequest.New().Post(url).Timeout(time.Second*10).Retry(5, 5*time.Second).SendStruct(&data).End()
	if len(errs) > 0 {
		log.Error("SendLarkTextNotifyAtAll req err:", errs)
	} else {
		log.Info("SendLarkTextNotifyAtAll req:", body)
	}
}
