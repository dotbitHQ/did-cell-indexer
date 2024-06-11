package http_svr

import (
	"context"
	"did-cell-indexer/http_svr/handle"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	log = logger.NewLoggerDefault("http_svr", logger.LevelDebug, nil)
)

type HttpSvr struct {
	Ctx context.Context
	H   *handle.HttpHandle

	Address string
	engine  *gin.Engine
	srv     *http.Server
}

func (h *HttpSvr) Run() {
	h.engine = gin.New()
	h.initRouter()

	h.srv = &http.Server{
		Addr:    h.Address,
		Handler: h.engine,
	}
	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			log.Error("ListenAndServe err:", err)
		}
	}()
}

func (h *HttpSvr) Shutdown() {
	if h.srv != nil {
		log.Warn("HttpSvr Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("Shutdown 1 err:", err.Error())
		}
	}
}
