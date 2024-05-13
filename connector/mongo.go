package connector

import (
	"context"
	"errors"
	"fmt"

	"go-micro.dev/v4/config"
	"go-micro.dev/v4/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo : Global Database Connector
var Mongo = new(mongoConnector)

type mongoConnector struct {
	client       *mongo.Client
	databaseName string // 业务库名 (认证库用admin,此处为配置文件中的业务库名)
}

type mongoConfiguration struct {
	AuthDatabase string `json:"auth_database"`
	Database     string `json:"database"`
	Address      string `json:"address"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func (c *mongoConnector) init(option *mongoConfiguration) error {
	var err error
	c.databaseName = option.Database
	if option.AuthDatabase == "" {
		option.AuthDatabase = "admin"
	}
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		option.Username, option.Password, option.Address, option.Port,
		option.AuthDatabase)
	if c.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn)); err != nil {
		return fmt.Errorf("failed to connect mongo database: %s", dsn)
	}
	return c.client.Ping(context.TODO(), nil)
}

func (c *mongoConnector) Instance() *mongo.Database {
	return c.client.Database(c.databaseName)
}

func (c *mongoConnector) Disconnect() error {
	return c.client.Disconnect(context.TODO())
}

func (c *mongoConnector) reInit() error {
	var option mongoConfiguration
	if err := config.Get("mongo").Scan(&option); err != nil {
		return errors.New("not found mongo configuration")
	}
	if err := c.init(&option); err != nil {
		return fmt.Errorf("init mongo connector error: %s", err.Error())
	}

	logger.Infof("mongo connect to: %v", option)

	return nil
}

func (c *mongoConnector) Initializer() error {
	// w, _ := config.Watch("mongo")
	// run.Async(context.Background(), func() {
	// 	for {
	// 		if _, err := w.Next(); err != nil {
	// 			logger.Errorf("watch mongo error: %v", err)
	// 		}
	// 		time.Sleep(time.Millisecond * 100)
	// 		if err := c.Disconnect(); err != nil {
	// 			logger.Infof("mongo close error:", err.Error())
	// 		}
	// 		if err := c.reInit(); err != nil {
	// 			logger.Errorf(err.Error())
	// 		}
	// 	}
	// })
	return c.reInit()
}

// IsMongoErrNilValue 值为空
func IsMongoErrNilValue(err error) bool {
	return errors.Is(err, mongo.ErrNilValue)
}

// IsMongoErrNilDocument document 为空
func IsMongoErrNilDocument(err error) bool {
	return errors.Is(err, mongo.ErrNilDocument)
}

// IsMongoErrEmptySlice : ...
func IsMongoErrEmptySlice(err error) bool {
	return errors.Is(err, mongo.ErrEmptySlice)
}

// IsMongoErrNilCursor : ...
func IsMongoErrNilCursor(err error) bool {
	return errors.Is(err, mongo.ErrNilCursor)
}

// IsMongoErrNoDocuments : ...
func IsMongoErrNoDocuments(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

// IsMongoErrWrongClient : ...
func IsMongoErrWrongClient(err error) bool {
	return errors.Is(err, mongo.ErrWrongClient)
}

// IsMongoErrClientDisconnected : ...
func IsMongoErrClientDisconnected(err error) bool {
	return errors.Is(err, mongo.ErrClientDisconnected)
}

// IsMongoErrInvalidIndexValue : ...
func IsMongoErrInvalidIndexValue(err error) bool {
	return errors.Is(err, mongo.ErrInvalidIndexValue)
}
