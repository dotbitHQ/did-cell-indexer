package main

import (
	"context"
	"did-cell-indexer/block_parser"
	"did-cell-indexer/cache"
	"did-cell-indexer/config"
	"did-cell-indexer/dao"
	"did-cell-indexer/http_svr"
	"did-cell-indexer/http_svr/handle"
	"did-cell-indexer/prometheus"

	"fmt"

	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/scorpiotzh/toolib"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
	"time"
)

var (
	log               = logger.NewLogger("main", logger.LevelDebug)
	exit              = make(chan struct{})
	ctxServer, cancel = context.WithCancel(context.Background())
	wgServer          = sync.WaitGroup{}
)

func main() {
	log.Debugf("server start：")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {
	// config file
	configFilePath := ctx.String("config")
	log.Info("configFilePath:", configFilePath)
	if err := config.InitCfg(configFilePath); err != nil {
		return err
	}

	// config file watcher
	watcher, err := config.AddCfgFileWatcher(configFilePath)
	if err != nil {
		return err
	}
	// ============= service start =============

	// sentry
	if err := http_api.SentryInit(config.Cfg.Notify.SentryDsn); err != nil {
		return fmt.Errorf("SentryInit err: %s", err.Error())
	}
	defer http_api.RecoverPanic()

	// prometheus
	prometheus.Init()
	prometheus.Tools.Run()

	// db
	log.Info("config.Cfg.DB.Mysql:", config.Cfg.DB.Mysql)
	dbDao, err := dao.NewGormDB(config.Cfg.DB.Mysql)
	if err != nil {
		return fmt.Errorf("dao.NewGormDB err: %s", err.Error())
	}

	// das core
	dasCore, dasCache, err := config.InitDasCore(ctxServer, &wgServer)
	if err != nil {
		return fmt.Errorf("config.InitDasCore err: %s", err.Error())
	}

	txBuilderBase, serverScript, err := config.InitTxBuilder(ctxServer, dasCore)
	if err != nil {
		return fmt.Errorf("config.InitDasTxBuilderBase err: %s", err.Error())
	}
	// block paser
	bp := block_parser.BlockParser{
		DasCore:            dasCore,
		CurrentBlockNumber: config.Cfg.Chain.Ckb.CurrentBlockNumber,
		DbDao:              dbDao,
		ConcurrencyNum:     config.Cfg.Chain.Ckb.ConcurrencyNum,
		ConfirmNum:         config.Cfg.Chain.Ckb.ConfirmNum,
		Ctx:                ctxServer,
		Wg:                 &wgServer,
	}
	if err := bp.Run(); err != nil {
		return fmt.Errorf("block_parser run err: %s", err.Error())
	}

	// redis
	red, err := toolib.NewRedisClient(config.Cfg.DB.Redis.Addr, config.Cfg.DB.Redis.Password, config.Cfg.DB.Redis.DbNum)
	if err != nil {
		return fmt.Errorf("NewRedisClient err:%s", err.Error())
	} else {
		log.Info("redis ok")
	}
	rc := &cache.RedisCache{
		Ctx: ctxServer,
		Red: red,
	}

	// http
	httpSvr := http_svr.HttpSvr{
		Ctx:     ctxServer,
		Address: config.Cfg.Server.HttpPort,
		H: &handle.HttpHandle{
			Ctx:           ctxServer,
			DbDao:         dbDao,
			RC:            rc,
			DasCore:       dasCore,
			DasCache:      dasCache,
			ServerScript:  serverScript,
			TxBuilderBase: txBuilderBase,
		},
	}
	httpSvr.Run()

	// ============= service end =============
	toolib.ExitMonitoring(func(sig os.Signal) {
		log.Warn("ExitMonitoring:", sig.String())
		if watcher != nil {
			log.Warn("close watcher ... ")
			_ = watcher.Close()
		}
		httpSvr.Shutdown()
		cancel()
		wgServer.Wait()

		log.Warn("success exit server. bye bye!")
		time.Sleep(time.Second)
		exit <- struct{}{}
	})

	<-exit
	return nil
}
