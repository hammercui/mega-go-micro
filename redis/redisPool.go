/**
redis连接池
 */
package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisPool struct {
	key string
	addr string
	dbIndex int
	pools []*redis.Client
}

func NewRedisPool(key string,addr string,dbIndex int,client *redis.Client) *RedisPool  {
	return &RedisPool{
		key:     key,
		addr:    addr,
		dbIndex: dbIndex,
		pools:   []*redis.Client{client},
	}
}

func (p *RedisPool) GetClient() *redis.Client  {
	if len(p.pools) <1{
		panic(fmt.Sprintf("redis connect is not exist,addr:%s",p.addr))
	}
	return p.pools[0]
}

func (p *RedisPool) PushClient(client *redis.Client)  {
	p.pools = append(p.pools,client)
}
