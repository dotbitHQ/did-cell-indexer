package http_svr

import (
	"context"
	"did-cell-indexer/http_svr/handle"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	log = logger.NewLogger("http_svr", logger.LevelDebug)
)

type HttpSvr struct {
	Ctx context.Context
	H   *handle.HttpHandle

	Address string
	engine  *gin.Engine
	srv     *http.Server

	InternalAddress string
	internalEngine  *gin.Engine
	internalSrv     *http.Server
}

func (h *HttpSvr) Run() {
	h.engine = gin.New()
	h.internalEngine = gin.New()
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

	h.internalSrv = &http.Server{
		Addr:    h.InternalAddress,
		Handler: h.internalEngine,
	}
	go func() {
		if err := h.internalSrv.ListenAndServe(); err != nil {
			log.Error("http_server internal run err:", err)
		}
	}()

}

func (h *HttpSvr) Shutdown() {
	if h.srv != nil {
		log.Warn("HttpSvr Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("Shutdown 1 err:", err.Error())
		}
		if err := h.internalSrv.Shutdown(h.Ctx); err != nil {
			log.Error("Shutdown 2 err:", err.Error())
		}
	}
}
