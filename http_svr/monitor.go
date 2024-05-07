package http_svr

import (
	"bytes"
	"did-cell-indexer/prometheus"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func DoMonitorLog(method string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()
		statusCode := ctx.Writer.Status()

		var resp http_api.ApiResp
		if statusCode == http.StatusOK && blw.body.String() != "" {
			if err := json.Unmarshal(blw.body.Bytes(), &resp); err != nil {
				log.Warn("DoMonitorLog Unmarshal err:", method, err)
				return
			}
			if resp.ErrNo != http_api.ApiCodeSuccess {
				log.Warn("DoMonitorLog:", method, resp.ErrNo, resp.ErrMsg)
			}
		}
		prometheus.Tools.Metrics.Api().WithLabelValues(method, fmt.Sprint(statusCode), fmt.Sprint(resp.ErrNo), resp.ErrMsg).Observe(time.Since(startTime).Seconds())
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyWriter) Write(bys []byte) (int, error) {
	b.body.Write(bys)
	return b.ResponseWriter.Write(bys)
}
