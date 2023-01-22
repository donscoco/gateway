package redis

import (
	"github.com/go-redis/redis"
	"time"
)

// todo 暂时先做单节点的临时demo。后续修改兼容集群
type ReidsManager struct {
	Addrs     []string
	client    *redis.Client
	isCluster bool
}

func NewRedisManager(addrs []string) (rm *ReidsManager) {
	// todo check
	opt := &redis.Options{
		Addr:               addrs[0],
		DB:                 0,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
	rm = new(ReidsManager)
	rm.Addrs = addrs
	rm.client = redis.NewClient(opt)
	rm.client.IncrBy()
	return

}
