/*
 * Copyright (C) Zenkoo, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file runtime.go
 * @package runtime
 * @author Dr.NP <zhanghao@liangyu.ltd>
 * @since 05/07/2024
 */

package runtime

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/alexlast/bunzap"
	brkKafka "github.com/go-micro/plugins/v4/broker/kafka"
	brkNats "github.com/go-micro/plugins/v4/broker/nats"
	brkRabbitmq "github.com/go-micro/plugins/v4/broker/rabbitmq"
	chRedis "github.com/go-micro/plugins/v4/cache/redis"
	srcConsul "github.com/go-micro/plugins/v4/config/source/consul"
	rgConsul "github.com/go-micro/plugins/v4/registry/consul"
	rgEtcd "github.com/go-micro/plugins/v4/registry/etcd"
	stConsul "github.com/go-micro/plugins/v4/store/consul"
	stRedis "github.com/go-micro/plugins/v4/store/redis"
	redis8 "github.com/go-redis/redis/v8"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/zenkoo-live/svc.base/middleware/session"
	"github.com/zenkoo-live/svc.base/zlogger"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/config"
	srcFile "go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	rg           registry.Registry
	brk          broker.Broker
	ch           cache.Cache
	st           store.Store
	db           *bun.DB
	mdb          *mongo.Client
	rdb          *redis.Client
	fb           *fiber.App
	fbAddress    string
	zaplogger    *zlogger.Zaplog
	env          string
	errorHandler fiber.ErrorHandler
)

func Init() error {
	var (
		err, errs error
	)

	// Load config
	cfgFilePath := config.Get("config", "file", "path").String("")
	cfgConsulPrefix := config.Get("config", "consul", "prefix").String("")

	// From local file
	if cfgFilePath != "" {
		cfgFile := srcFile.NewSource(srcFile.WithPath(cfgFilePath))
		err = config.Load(cfgFile)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	// From remote consul
	if cfgConsulPrefix != "" {
		cfgConsulAddress := config.Get("config", "consul", "address").String("")
		cfgConsul := srcConsul.NewSource(
			srcConsul.WithAddress(cfgConsulAddress),
			srcConsul.WithPrefix(cfgConsulPrefix),
			srcConsul.StripPrefix(true),
		)
		err = config.Load(cfgConsul)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	cfg := &Config{}
	config.Get(BaseConfigKey).Scan(cfg)

	// Env
	env = strings.ToLower(config.Get("env").String(""))

	// Logger
	loggerOpts := []logger.Option{
		zlogger.WithCallerSkip(2),
	}
	if cfg.Logger != nil {
		if cfg.Logger.Debug {
			loggerOpts = append(loggerOpts, logger.WithLevel(logger.DebugLevel))
		}

		if cfg.Logger.Silence {
			loggerOpts = append(loggerOpts, zlogger.WithEcho(false))
		}

		if cfg.Logger.DiscardDisk {
			loggerOpts = append(loggerOpts, zlogger.WithDisk(false))
		}

		if cfg.Logger.BasePath != "" {
			loggerOpts = append(loggerOpts, zlogger.WithBasePath(cfg.Logger.BasePath))
		}
	}

	zaplogger, _ = zlogger.NewLogger(loggerOpts...)
	logger.DefaultLogger = zaplogger

	// Registry
	if cfg.Registry != nil {
		rg, err = initRegistry(cfg.Registry)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			registry.DefaultRegistry = rg
		}
	} else {
		rg = registry.DefaultRegistry
	}

	// Broker
	if cfg.Broker != nil {
		brk, err = initBroker(cfg.Broker)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			broker.DefaultBroker = brk
		}
	} else {
		brk = broker.DefaultBroker
	}

	// Cache
	if cfg.Cache != nil {
		ch, err = initCache(cfg.Cache)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			cache.DefaultCache = ch
		}
	} else {
		ch = cache.DefaultCache
	}

	// Store
	if cfg.Store != nil {
		st, err = initStore(cfg.Store)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			store.DefaultStore = st
		}
	} else {
		st = store.DefaultStore
	}

	// Database
	if cfg.Database != nil {
		db, err = initDatabase(cfg.Database)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	// Mongo
	if cfg.Mongo != nil {
		mdb, err = initMongo(cfg.Mongo)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	// Redis
	if cfg.Redis != nil {
		rdb, err = initRedis(cfg.Redis)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	// Fiber
	if cfg.Fiber != nil {
		fb, err = initFiber(cfg.Fiber)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	// Session
	if cfg.Session != nil {
		err = initSession(cfg.Session)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func initRegistry(cfg *configRegistry) (registry.Registry, error) {
	if cfg == nil {
		return nil, errors.New("empty registry configuration")
	}

	var trg registry.Registry

	switch strings.ToLower(cfg.Driver) {
	case "etcd":
		// ETCD v3
		trg = rgEtcd.NewRegistry(
			registry.Addrs(cfg.Address...),
		)
	default:
		// Consul
		trg = rgConsul.NewRegistry(
			registry.Addrs(cfg.Address...),
		)
	}

	logger.Infof("registry <%s> initialized", cfg.Driver)

	return trg, nil
}

func initBroker(cfg *configBroker) (broker.Broker, error) {
	if cfg == nil {
		return nil, errors.New("empty broker configuration")
	}

	var tbrk broker.Broker

	switch strings.ToLower(cfg.Driver) {
	case "kafka":
		// Kafka
		tbrk = brkKafka.NewBroker(
			broker.Addrs(cfg.Address...),
		)
	case "rabbitmq":
		tbrk = brkRabbitmq.NewBroker(
			broker.Addrs(cfg.Address...),
		)
	default:
		// Nats
		tbrk = brkNats.NewBroker(
			broker.Addrs(cfg.Address...),
		)
	}

	if err := tbrk.Init(); err != nil {
		logger.Fatalf("Broker Init error: %v", err)

		return nil, err
	}

	if err := tbrk.Connect(); err != nil {
		logger.Fatalf("Broker Connect error: %v", err)

		return nil, err
	}

	logger.Infof("broker <%s> initialized", cfg.Driver)

	return tbrk, nil
}

func initCache(cfg *configCache) (cache.Cache, error) {
	if cfg == nil {
		return nil, errors.New("empty cache configuration")
	}

	var tch cache.Cache

	switch strings.ToLower(cfg.Driver) {
	case "memory":
	case "redis":
		// Redis
		tch = chRedis.NewCache(
			chRedis.WithRedisOptions(
				redis8.UniversalOptions{
					Addrs:    cfg.Address,
					DB:       cfg.DB,
					Password: cfg.Password,
				},
			),
		)
	default:
		tch = cache.DefaultCache
	}

	logger.Infof("cache <%s> initialized", cfg.Driver)

	return tch, nil
}

func initStore(cfg *configStore) (store.Store, error) {
	if cfg == nil {
		return nil, errors.New("empty store configuration")
	}

	var tst store.Store

	switch strings.ToLower(cfg.Driver) {
	case "consul":
		// Consul
		// TODO: DO NOT USE ME
		tst = stConsul.NewStore()
	case "redis":
		// Redis
		tst = stRedis.NewStore(
			stRedis.WithRedisOptions(
				redis8.UniversalOptions{
					Addrs:    cfg.Address,
					DB:       cfg.DB,
					Password: cfg.Password,
				},
			),
		)
	default:
		tst = store.DefaultStore
	}

	logger.Infof("store <%s> initialized", cfg.Driver)

	return tst, nil
}

func initDatabase(cfg *configDatabase) (*bun.DB, error) {
	if cfg == nil {
		return nil, errors.New("empty database configuration")
	}

	var (
		sqldb *sql.DB
		err   error
		tdb   *bun.DB
	)

	switch strings.ToLower(cfg.Driver) {
	case "mysql":
		// MySQL
		sqldb, err = sql.Open("mysql", cfg.DSN)
		if err != nil {
			return nil, err
		}

		tdb = bun.NewDB(sqldb, mysqldialect.New())
	case "mssql":
		// MS-SQLServer
		sqldb, err = sql.Open("sqlserver", cfg.DSN)
		if err != nil {
			return nil, err
		}

		tdb = bun.NewDB(sqldb, mssqldialect.New())
	case "sqlite":
		// SQLite
		sqldb, err = sql.Open("sqlite3", cfg.DSN)
		if err != nil {
			return nil, err
		}

		tdb = bun.NewDB(sqldb, sqlitedialect.New())
	default:
		// PostgreSQL
		sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
		tdb = bun.NewDB(sqldb, pgdialect.New())
	}

	err = tdb.Ping()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		tdb.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
			Logger: zaplogger.Zap(),
		}))
	} else if cfg.SlowQueryDuration > 0 {
		tdb.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
			Logger:       zaplogger.Zap(),
			SlowDuration: time.Duration(cfg.SlowQueryDuration) * time.Millisecond,
		}))
	}

	logger.Infof("database <%s> initialized", cfg.Driver)

	return tdb, nil
}

func initMongo(cfg *configMongo) (*mongo.Client, error) {
	if cfg == nil {
		return nil, errors.New("empty mongo configuration")
	}

	var (
		err error
		tdb *mongo.Client
	)

	tdb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DSN))
	if err != nil {
		return nil, err
	}

	logger.Infof("mongodb initialized")

	return tdb, nil
}

func initRedis(cfg *configRedis) (*redis.Client, error) {
	if cfg == nil {
		return nil, errors.New("empty redis configuration")
	}

	var (
		err error
		tdb *redis.Client
	)

	tdb = redis.NewClient(&redis.Options{
		Addr:       cfg.Address,
		Password:   cfg.Password,
		DB:         cfg.DB,
		PoolSize:   cfg.PoolSize,
		MaxRetries: cfg.MaxRetries,
	})
	_, err = tdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	logger.Infof("redis initialized")

	return tdb, nil
}

func initFiber(cfg *configFiber) (*fiber.App, error) {
	if cfg == nil {
		return nil, errors.New("empty fiber configuration")
	}

	fc := fiber.Config{
		Prefork:               false,
		DisableStartupMessage: true,
		Network:               "tcp",
		DisableKeepalive:      cfg.DisableKeepAlive,
		StrictRouting:         cfg.StrictRouting,
		CaseSensitive:         cfg.CaseSensitive,
		ETag:                  cfg.Etag,
		BodyLimit:             cfg.BodyLimit,
		Concurrency:           cfg.Concurrency,
		ReadBufferSize:        cfg.ReadBufferSize,
		WriteBufferSize:       cfg.WriteBufferSize,
		ErrorHandler:          errorHandler,
	}
	tfb := fiber.New(fc)
	if cfg.EnableStackTrace {
		tfb.Use(recover.New(
			recover.Config{
				EnableStackTrace: true,
			},
		))
	} else {
		tfb.Use(recover.New(
			recover.ConfigDefault,
		))
	}

	tfb.Use(
		cors.New(),
		favicon.New(),
		requestid.New(),
		fiberzap.New(
			fiberzap.Config{
				Logger: zaplogger.Zap(),
				Fields: []string{"latency", "status", "method", "url", "ip", "requestId", "error"},
			},
		),
	)

	if cfg.EnableSwagger {
		tfb.All("/docs/*", swagger.New(swagger.ConfigDefault))
	}

	if cfg.Address != "" {
		fbAddress = cfg.Address
	} else {
		fbAddress = DefaultHTTPAdvertiseAddr
	}

	return tfb, nil
}

func initSession(cfg *configSession) error {
	if cfg == nil {
		return errors.New("empty session configuration")
	}

	if fb != nil {
		r := rdb
		if r == nil {
			// TODO: create new redis client from config
			r = redis.NewClient(&redis.Options{})
		}

		mw := session.New(
			session.Config{
				IDSource:   cfg.IDSource,
				IDKey:      cfg.IDKey,
				IDPrefix:   cfg.IDPrefix,
				Expiration: time.Second * time.Duration(cfg.Expiration),
				Storage:    r,
			},
		)
		fb.Use(mw)
	}

	logger.Info("fiber session initialized")

	return nil
}

func Registry() registry.Registry {
	return rg
}

func Broker() broker.Broker {
	return brk
}

func Cache() cache.Cache {
	return ch
}

func Store() store.Store {
	return st
}

func DB() *bun.DB {
	return db
}

func Mongo() *mongo.Client {
	return mdb
}

func Redis() *redis.Client {
	return rdb
}

func Fiber() *fiber.App {
	return fb
}

func Logger() logger.Logger {
	return zaplogger
}

func Zap() *zap.Logger {
	return zaplogger.Zap()
}

func Env() string {
	return env
}

func AppendEnv() string {
	if env != "" {
		return "::" + env
	}

	return ""
}

func StartHTTP() {
	if fb != nil {
		go func() {
			logger.Infof("fiber listening on %s", fbAddress)
			err := fb.Listen(fbAddress)
			if err != nil {
				logger.Errorf("fiber listen failed on %s : %s", fbAddress, err.Error())

				return
			}

			logger.Infof("fiber closed")
		}()
	}
}

func StopHTTP() {
	if fb != nil {
		logger.Info("stopping fiber server")
		fb.Shutdown()
	}
}

func SetFiberErrorHandler(handler fiber.ErrorHandler) {
	errorHandler = handler
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
