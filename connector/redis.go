package connector

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/logger"
)

// Redis : Global Redis Connector
var Redis = new(redisConnector)

type redisConnector struct {
	client *redis.Client
}

type redisConfiguration struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	DB         string `json:"db"`
	Password   string `json:"password"`
	PoolSize   string `json:"pool_size"`
	MaxRetries string `json:"max_retries"`
}

func (c *redisConnector) init(option *redisConfiguration) error {
	db, _ := strconv.Atoi(option.DB)
	poolSize, _ := strconv.Atoi(option.PoolSize)
	maxRetries, _ := strconv.Atoi(option.MaxRetries)
	if poolSize == 0 {
		poolSize = 5
	}
	if maxRetries == 0 {
		maxRetries = 3
	}
	c.client = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%v:%v", option.Host, option.Port),
		Password:   option.Password,
		DB:         db,
		PoolSize:   poolSize,
		MaxRetries: maxRetries,
	})
	if _, err := c.client.Ping(context.TODO()).Result(); err != nil {
		return fmt.Errorf("redis init ping err:%v", err)
	}
	return nil
}

func (c *redisConnector) Instance() *redis.Client {
	return c.client
}

func (c *redisConnector) reInit() error {
	var option redisConfiguration
	if err := config.Get("redis").Scan(&option); err != nil {
		return fmt.Errorf("not found redis configuration: %s", err.Error())
	}
	if err := c.init(&option); err != nil {
		return fmt.Errorf("init redis connector error: %s", err.Error())
	}
	logger.Infof("redis connect to: %v", option)

	return nil
}

func (c *redisConnector) Initializer() error {
	// w, _ := config.Watch("redis")
	// run.Async(context.Background(), func() {
	// 	for {
	// 		if _, err := w.Next(); err != nil {
	// 			logger.Errorf("watch redis error: %v", err)
	// 		}
	// 		time.Sleep(time.Millisecond * 100)
	// 		if err := c.client.Close(); err != nil {
	// 			logger.Infof("redis close error:", err.Error())
	// 		}
	// 		if err := c.reInit(); err != nil {
	// 			logger.Errorf(err.Error())
	// 		}
	// 	}
	// })
	return c.reInit()
}

// IsNil  Nil reply returned by Redis when key does not exist.
func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
