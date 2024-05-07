package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"time"
)

const (
	lockTime    = 180
	lockTicker  = 10
	lockOrderId = "lock:order_id:"
	lock        = "lock:"
)

var ErrDistributedLockPreemption = errors.New("distributed lock preemption")

func (r *RedisCache) LockWithRedis(orderId string) error {
	log.Info("LockWithRedis:", orderId)
	ret := r.Red.SetNX(lockOrderId+orderId, orderId, time.Second*lockTime)
	if err := ret.Err(); err != nil {
		return fmt.Errorf("redis set order nx-->%s", err.Error())
	}
	if !ret.Val() {
		return ErrDistributedLockPreemption
	}
	return nil
}

func (r *RedisCache) UnLockWithRedis(orderId string) error {
	log.Info("UnLockWithRedis:", orderId)
	ret := r.Red.Del(lockOrderId + orderId)
	if err := ret.Err(); err != nil {
		return fmt.Errorf("redis del order nx-->%s", err.Error())
	}
	return nil
}

func (r *RedisCache) DoLockExpire(ctx context.Context, orderId string) {
	ticker := time.NewTicker(time.Second * lockTicker)
	count := 0
	go func() {
		defer http_api.RecoverPanic()
		for {
			select {
			case <-ticker.C:
				ok, err := r.Red.Expire(lockOrderId+orderId, time.Second*lockTime).Result()
				if err != nil {
					log.Error("DoLockExpire err: ", err.Error(), orderId)
				} else if ok {
					count++
				}
				log.Infof("DoLockExpire: %s %d %p", orderId, count, &count)
			case <-ctx.Done():
				log.Info("DoLockExpire done:", orderId)
				return
			}
		}
	}()

}
