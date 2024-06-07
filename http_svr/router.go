package http_svr

import (
	"did-cell-indexer/config"
	"encoding/json"
	"github.com/dotbitHQ/das-lib/http_api"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

func (h *HttpSvr) initRouter() {
	if len(config.Cfg.Origins) > 0 {
		toolib.AllowOriginList = append(toolib.AllowOriginList, config.Cfg.Origins...)
	}
	h.engine.Use(toolib.MiddlewareCors())
	h.engine.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	v1 := h.engine.Group("v1")
	{
		// cache
		shortExpireTime, shortDataTime, lockTime := time.Second*5, time.Minute*3, time.Minute
		cacheHandleShort := toolib.MiddlewareCacheByRedis(h.H.RC.Red, false, shortDataTime, lockTime, shortExpireTime, respHandle)
		//longExpireTime, longDataTime := time.Second*15, time.Minute*10
		//cacheHandleLong := toolib.MiddlewareCacheByRedis(h.rc.GetRedisClient(), false, longDataTime, lockTime, longExpireTime, respHandle)
		//cacheHandleShortCookies := toolib.MiddlewareCacheByRedis(h.rc.GetRedisClient(), true, shortDataTime, lockTime, shortExpireTime, respHandle)
		v1.POST("/account/list", DoMonitorLog("account_list"), cacheHandleShort, h.H.AccountList) //
		v1.POST("/record/list", DoMonitorLog("record_list"), cacheHandleShort, h.H.RecordList)    //
		// operate
		v1.POST("/edit/owner", DoMonitorLog("edit_owner"), h.H.EditOwner)
		v1.POST("/edit/record", DoMonitorLog("edit_record"), h.H.EditRecord)
		v1.POST("/recycle", DoMonitorLog("recycle"), h.H.Recycle)
		v1.POST("/tx/send", DoMonitorLog("tx_send"), h.H.TxSend)

	}

}

func respHandle(c *gin.Context, res string, err error) {
	if err != nil {
		log.Error("respHandle err:", err.Error())
		c.AbortWithStatusJSON(http.StatusOK, http_api.ApiRespErr(http.StatusInternalServerError, err.Error()))
	} else if res != "" {
		var respMap map[string]interface{}
		_ = json.Unmarshal([]byte(res), &respMap)
		c.AbortWithStatusJSON(http.StatusOK, respMap)
	}
}
