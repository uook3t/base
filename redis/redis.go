package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCli struct {
	cli *redis.Client
}

var rdb *RedisCli

func InitRedis(addr string) {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
		// Password: conf.Password, // 密码
		// DB:       0,             // 数据库
		PoolSize: 10, // 连接池大小
	})
	rdb = &RedisCli{cli: cli}
}

func GetRedis() *RedisCli {
	return rdb
}

func (c *RedisCli) Get(ctx context.Context, k string) (string, error) {
	cmd := c.cli.Get(ctx, k)
	return cmd.Val(), cmd.Err()
}

func (c *RedisCli) Set(ctx context.Context, k, v string, expire time.Duration) error {
	cmd := c.cli.Set(ctx, k, v, expire)
	return cmd.Err()
}

func (c *RedisCli) MSet(ctx context.Context, kvs map[string]string, expires map[string]time.Duration) error {
	pipeline := c.cli.Pipeline()
	for k, v := range kvs {
		expire, ok := expires[k]
		if !ok {
			expire = defaultExpireTime
		}
		pipeline.Set(ctx, k, v, expire)
	}
	cmds, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func (c *RedisCli) SetNX(ctx context.Context, k, v string, expire time.Duration) (bool, error) {
	cmd := c.cli.SetNX(ctx, k, v, expire)
	return cmd.Val(), cmd.Err()
}

func (c *RedisCli) Del(ctx context.Context, ks ...string) error {
	cmd := c.cli.Del(ctx, ks...)
	return cmd.Err()
}
