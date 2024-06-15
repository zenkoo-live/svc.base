/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file config.go
 * @package session
 * @author Dr.NP <conan.np@gmail.com>
 * @since 06/14/2024
 */

package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool
	// Session ID source
	IDSource string
	// Session ID key
	IDKey string
	// Session ID prefix
	IDPrefix string
	// Session Data key
	DataKey string
	// Expiration duration
	Expiration time.Duration
	// Auto refresh expiry
	AutoRefreshExpiration bool
	// Strict session data
	StrictAuth string
	// Storage, Redis client
	Storage *redis.Client
}

var ConfigDefault = Config{
	Next:                  nil,
	IDSource:              "cookie",
	IDKey:                 "session_id",
	IDPrefix:              "sess_",
	DataKey:               "session_data",
	Expiration:            time.Second * 3600,
	StrictAuth:            "",
	AutoRefreshExpiration: false,
	Storage:               redis.NewClient(&redis.Options{}),
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.IDSource == "" {
		cfg.IDSource = ConfigDefault.IDSource
	}

	if cfg.IDKey == "" {
		cfg.IDKey = ConfigDefault.IDKey
	}

	if cfg.IDPrefix == "" {
		cfg.IDPrefix = ConfigDefault.IDPrefix
	}

	if cfg.DataKey == "" {
		cfg.DataKey = ConfigDefault.DataKey
	}

	if cfg.Expiration == 0 {
		cfg.Expiration = ConfigDefault.Expiration
	}

	if cfg.Storage == nil {
		cfg.Storage = ConfigDefault.Storage
	}

	return cfg
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
