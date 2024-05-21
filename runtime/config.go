/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file config.go
 * @package runtime
 * @author Dr.NP <zhanghao@liangyu.ltd>
 * @since 05/07/2024
 */

package runtime

import (
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	BaseConfigKey  = "base"
	DefaultAppName = "zenkoo"
	DefaultSvcName = "svc"
	DefaultVersion = "latest"

	DefaultRegistryDriver = "consul"
	DefaultBrokerDriver   = "nats"
	DefaultCacheDriver    = "redis"

	DefaultHTTPAdvertiseAddr = ":9990"
	DefaultGRPCAdvertiseAddr = ":9991"
)

type configRegistry struct {
	Driver  string   `json:"driver" mapstructure:"driver"`
	Address []string `json:"address" mapstructure:"address"`
}

type configBroker struct {
	Driver  string   `json:"driver" mapstructure:"driver"`
	Address []string `json:"address" mapstructure:"address"`
}

type configCache struct {
	Driver  string `json:"driver" mapstructure:"driver"`
	Address string `json:"address" mapstructure:"address"`
}

type configDatabase struct {
	Driver            string `json:"driver" mapstructure:"driver"`
	DSN               string `json:"dsn" mapstructure:"dsn"`
	Debug             bool   `json:"debug" mapstructure:"debug"`
	SlowQueryDuration int    `json:"slow_query_duration" mapstructure:"slow_query_duration"`
}

type configFiber struct {
	Address          string `json:"address" mapstructure:"address"`
	StrictRouting    bool   `json:"strict_routing" mapstructure:"strict_routing"`
	CaseSensitive    bool   `json:"case_sensitive" mapstructure:"case_sensitive"`
	Etag             bool   `json:"etag" mapstructure:"etag"`
	BodyLimit        int    `json:"body_limit" mapstructure:"body_limit"`
	Concurrency      int    `json:"concurrency" mapstructure:"concurrency"`
	ReadBufferSize   int    `json:"read_buffer_size" mapstructure:"read_buffer_size"`
	WriteBufferSize  int    `json:"write_buffer_size" mapstructure:"write_buffer_size"`
	DisableKeepAlive bool   `json:"disable_keep_alive" mapstructure:"disable_keep_alive"`
	EnableSwagger    bool   `json:"enable_swagger" mapstructure:"enable_swagger"`
	EnableStackTrace bool   `json:"enable_stack_trace" mapstructure:"enable_stack_trace"`
}

type configLogger struct {
	Debug       bool   `json:"debug" mapstructure:"debug"`
	Silence     bool   `json:"silence" mapstructure:"silence"`
	DiscardDisk bool   `json:"discard_disk" mapstructure:"discard_disk"`
	BasePath    string `json:"base_path" mapstructure:"base_path"`
}

type configMongo struct {
	DSN string `json:"dsn" mapstructure:"dsn"`
}

type configRedis struct {
	Address    string `json:"address" mapstructure:"address"`
	Password   string `json:"password" mapstructure:"password"`
	PoolSize   int    `json:"pool_size" mapstructure:"pool_size"`
	MaxRetries int    `json:"max_retries" mapstructure:"max_retries"`
	DB         int    `json:"db" mapstructure:"db"`
}

type Config struct {
	Registry *configRegistry `json:"registry" mapstructure:"registry"`
	Broker   *configBroker   `json:"broker" mapstructure:"broker"`
	Cache    *configCache    `json:"cache" mapstructure:"cache"`
	Database *configDatabase `json:"database" mapstructure:"database"`
	Mongo    *configMongo    `json:"mongo" mapstructure:"mongo"`
	Redis    *configRedis    `json:"redis" mapstructure:"redis"`
	Fiber    *configFiber    `json:"fiber" mapstructure:"fiber"`
	Logger   *configLogger   `json:"logger" mapstructure:"logger"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
