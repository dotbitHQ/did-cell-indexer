package cache

import (
	"crypto/md5"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"time"
)

func (r *RedisCache) getSignTxCacheKey(key string) string {
	return "sign:tx:" + key
}

func (r *RedisCache) GetSignTxCache(key string) (string, error) {
	if r.Red == nil {
		return "", fmt.Errorf("redis is nil")
	}
	key = r.getSignTxCacheKey(key)
	if txStr, err := r.Red.Get(key).Result(); err != nil {
		return "", err
	} else {
		return txStr, nil
	}
}

func (r *RedisCache) SetSignTxCache(key, txStr string) error {
	if r.Red == nil {
		return fmt.Errorf("redis is nil")
	}
	key = r.getSignTxCacheKey(key)
	if err := r.Red.Set(key, txStr, time.Minute*10).Err(); err != nil {
		return err
	}
	return nil
}

type SignInfoCache struct {
	Action    common.DidCellAction               `json:"action"`
	CkbAddr   string                             `json:"ckb_addr"`
	BuilderTx *txbuilder.DasTxBuilderTransaction `json:"builder_tx"`
}

func (s *SignInfoCache) SignKey() string {
	key := fmt.Sprintf("%s%s%d", s.CkbAddr, s.Action, time.Now().UnixNano())
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
