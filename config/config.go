package config

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/dotbitHQ/das-lib/remote_sign"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/fsnotify/fsnotify"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
	"sync"
	"time"
)

var (
	Cfg CfgServer
	log = logger.NewLogger("config", logger.LevelDebug)
)

func InitCfg(configFilePath string) error {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	log.Debug("config file path：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	log.Debug("config file：", toolib.JsonString(Cfg))
	return nil
}

func AddCfgFileWatcher(configFilePath string) (*fsnotify.Watcher, error) {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	return toolib.AddFileWatcher(configFilePath, func() {
		log.Debug("config file path：", configFilePath)
		if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
			log.Error("UnmarshalYamlFile err:", err.Error())
		}
		log.Debug("config file：", toolib.JsonString(Cfg))
	})
}

type CfgServer struct {
	Server struct {
		IsUpdate         bool              `json:"is_update" yaml:"is_update"`
		Name             string            `json:"name" yaml:"name"`
		Net              common.DasNetType `json:"net" yaml:"net"`
		HttpPort         string            `json:"http_port" yaml:"http_port"`
		PayServerAddress string            `json:"pay_server_address" yaml:"pay_server_address"`
		PayPrivate       string            `json:"pay_private" yaml:"pay_private"`
		RemoteSignApiUrl string            `json:"remote_sign_api_url" yaml:"remote_sign_api_url"`
		TxTeeRate        uint64            `json:"tx_fee_rate" yaml:"tx_fee_rate"`
	} `json:"server" yaml:"server"`
	Origins []string `json:"origins" yaml:"origins"`
	Notify  struct {
		SentryDsn             string `json:"sentry_dsn" yaml:"sentry_dsn"`
		LarkErrKey            string `json:"lark_err_key" yaml:"lark_err_key"`
		PrometheusPushGateway string `json:"prometheus_push_gateway" yaml:"prometheus_push_gateway"`
	} `json:"notify" yaml:"notify"`
	DB struct {
		Mysql DbMysql `json:"mysql" yaml:"mysql"`
		Redis struct {
			Addr     string `json:"addr" yaml:"addr"`
			Password string `json:"password" yaml:"password"`
			DbNum    int    `json:"db_num" yaml:"db_num"`
		} `json:"redis" yaml:"redis"`
	} `json:"db" yaml:"db"`
	Chain struct {
		Ckb struct {
			Node               string `json:"node" yaml:"node"`
			Addr               string `json:"addr" yaml:"addr"`
			CurrentBlockNumber uint64 `json:"current_block_number" yaml:"current_block_number"`
			ConfirmNum         uint64 `json:"confirm_num" yaml:"confirm_num"`
			ConcurrencyNum     uint64 `json:"concurrency_num" yaml:"concurrency_num"`
		} `json:"ckb" yaml:"ckb"`
	} `json:"chain" yaml:"chain"`
}

type DbMysql struct {
	Addr     string `json:"addr" yaml:"addr"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	DbName   string `json:"db_name" yaml:"db_name"`
}

func InitDasCore(ctx context.Context, wg *sync.WaitGroup) (*core.DasCore, *dascache.DasCache, error) {
	// ckb node
	ckbClient, err := rpc.DialWithIndexer(Cfg.Chain.Ckb.Node, Cfg.Chain.Ckb.Node)
	if err != nil {
		return nil, nil, fmt.Errorf("rpc.DialWithIndexer err: %s", err.Error())
	}
	log.Info("ckb node ok")

	// das init
	net := Cfg.Server.Net

	env := core.InitEnvOpt(net,
		common.DasContractNameConfigCellType,
		common.DasContractNameDispatchCellType,
		//common.DasContractNameBalanceCellType,
		common.DASContractNameEip712LibCellType,
		//common.DasContractNameDpCellType,
		//common.DasKeyListCellType,
		common.DasContractNameDidCellType,
		common.DasContractNameAlwaysSuccess,
	)
	ops := []core.DasCoreOption{
		core.WithClient(ckbClient),
		core.WithDasContractArgs(env.ContractArgs),
		core.WithDasContractCodeHash(env.ContractCodeHash),
		core.WithDasNetType(net),
		core.WithTHQCodeHash(env.THQCodeHash),
	}
	dasCore := core.NewDasCore(ctx, wg, ops...)
	dasCore.InitDasContract(env.MapContract)
	if err := dasCore.InitDasConfigCell(); err != nil {
		return nil, nil, fmt.Errorf("InitDasConfigCell err: %s", err.Error())
	}
	if err := dasCore.InitDasSoScript(); err != nil {
		return nil, nil, fmt.Errorf("InitDasSoScript err: %s", err.Error())
	}
	dasCore.RunAsyncDasContract(time.Minute * 3)   // contract outpoint
	dasCore.RunAsyncDasConfigCell(time.Minute * 5) // config cell outpoint
	dasCore.RunAsyncDasSoScript(time.Minute * 7)   // so

	log.Info("das contract ok")

	// das cache
	dasCache := dascache.NewDasCache(ctx, wg)
	dasCache.RunClearExpiredOutPoint(time.Minute * 15)
	log.Info("das cache ok")
	return dasCore, dasCache, nil
}

func InitTxBuilder(ctx context.Context, dasCore *core.DasCore) (*txbuilder.DasTxBuilderBase, *types.Script, error) {
	payServerAddressArgs := ""
	var serverScript *types.Script

	if Cfg.Server.PayServerAddress != "" {
		parseAddress, err := address.Parse(Cfg.Server.PayServerAddress)
		if err != nil {
			log.Error("pay server address.Parse err: ", err.Error())
		} else {
			payServerAddressArgs = common.Bytes2Hex(parseAddress.Script.Args)
			serverScript = parseAddress.Script
		}
	}
	var handleSign sign.HandleSignCkbMessage
	if Cfg.Server.RemoteSignApiUrl != "" && payServerAddressArgs != "" {
		//remoteSignClient, err := sign.NewClient(ctxServer, config.Cfg.Server.RemoteSignApiUrl)
		//if err != nil {
		//	return nil, nil, fmt.Errorf("sign.NewClient err: %s", err.Error())
		//}
		//handleSign = sign.RemoteSign(remoteSignClient, config.Cfg.Server.Net, payServerAddressArgs)
		handleSign = remote_sign.SignTxForCKBHandle(Cfg.Server.RemoteSignApiUrl, Cfg.Server.PayServerAddress)
	} else if Cfg.Server.PayPrivate != "" {
		handleSign = sign.LocalSign(Cfg.Server.PayPrivate)
	}
	txBuilderBase := txbuilder.NewDasTxBuilderBase(ctx, dasCore, handleSign, payServerAddressArgs)
	log.Info("tx builder ok")

	return txBuilderBase, serverScript, nil
}
