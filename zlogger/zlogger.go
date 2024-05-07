/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file zlogger.go
 * @package zlogger
 * @author Dr.NP <zhanghao@liangyu.ltd>
 * @since 05/08/2024
 */

package zlogger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go-micro.dev/v4/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaplog struct {
	cfg  zap.Config
	zap  *zap.Logger
	opts logger.Options
	sync.RWMutex
	fields map[string]interface{}
}

const (
	// TraceLevel : Custom level out of zap
	TraceLevel = -2
)

var (
	_defaultEncoderConfig = zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "caller",
		FunctionKey:   "function",
		StacktraceKey: "stack",
		LineEnding:    "",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime:       zapcore.RFC3339NanoTimeEncoder,
		EncodeTime:       customTimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "",
	}
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (l *zaplog) Init(opts ...logger.Option) error {
	var (
		err   error
		cores []zapcore.Core
	)

	for _, o := range opts {
		o(&l.opts)
	}

	skip, ok := l.opts.Context.Value(callerSkipKey{}).(int)
	if !ok || skip < 1 {
		skip = 1
	}

	basePath, ok := l.opts.Context.Value(basePathKey{}).(string)
	if !ok {
		basePath = "./logs"
	}

	echo, ok := l.opts.Context.Value(echoKey{}).(bool)
	if !ok {
		echo = true
	}

	disk, ok := l.opts.Context.Value(diskKey{}).(bool)
	if !ok {
		disk = true
	}

	// Levels
	showLevel := loggerToZapLevel(l.opts.Level)
	// TRACE
	lvlTrace := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= TraceLevel && lvl == TraceLevel
	})
	// DEBUG
	lvlDebug := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= zapcore.DebugLevel && lvl == zapcore.DebugLevel
	})
	// INFO
	lvlInfo := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= zapcore.InfoLevel && lvl == zapcore.InfoLevel
	})
	// WARN
	lvlWarn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= zapcore.WarnLevel && lvl == zapcore.WarnLevel
	})
	// ERROR
	lvlError := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= zapcore.ErrorLevel && lvl == zapcore.ErrorLevel
	})
	// FATAL
	lvlFatal := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return showLevel <= zapcore.FatalLevel && lvl >= zapcore.FatalLevel
	})

	// Syncers
	syncerStdout := zapcore.AddSync(os.Stdout)
	syncerStderr := zapcore.AddSync(os.Stderr)
	//syncerSilence := zapcore.AddSync(ioutil.Discard)

	// Disk syncers
	rollingDebug, err := NewRollingFile(basePath+"/debug", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingDebug := zapcore.AddSync(rollingDebug)
	rollingInfo, err := NewRollingFile(basePath+"/info", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingInfo := zapcore.AddSync(rollingInfo)
	rollingWarn, err := NewRollingFile(basePath+"/warn", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingWarn := zapcore.AddSync(rollingWarn)
	rollingError, err := NewRollingFile(basePath+"/error", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingError := zapcore.AddSync(rollingError)
	rollingFatal, err := NewRollingFile(basePath+"/fatal", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingFatal := zapcore.AddSync(rollingFatal)
	rollingTrace, err := NewRollingFile(basePath+"/trace", HourlyRolling)
	if err != nil {
		return err
	}

	syncerRollingTrace := zapcore.AddSync(rollingTrace)

	// Cores
	jsonEncoder := zapcore.NewJSONEncoder(_defaultEncoderConfig)
	consoleEncoder := zapcore.NewJSONEncoder(_defaultEncoderConfig)
	if disk {
		cores = append(
			cores,
			zapcore.NewCore(jsonEncoder, syncerRollingTrace, lvlTrace),
			zapcore.NewCore(jsonEncoder, syncerRollingDebug, lvlDebug),
			zapcore.NewCore(jsonEncoder, syncerRollingInfo, lvlInfo),
			zapcore.NewCore(jsonEncoder, syncerRollingWarn, lvlWarn),
			zapcore.NewCore(jsonEncoder, syncerRollingError, lvlError),
			zapcore.NewCore(jsonEncoder, syncerRollingFatal, lvlFatal),
			// Must say
			zapcore.NewCore(consoleEncoder, syncerStderr, lvlFatal),
		)
	}

	if echo {
		cores = append(
			cores,
			zapcore.NewCore(consoleEncoder, syncerStdout, lvlTrace),
			zapcore.NewCore(consoleEncoder, syncerStdout, lvlDebug),
			zapcore.NewCore(consoleEncoder, syncerStdout, lvlInfo),
			zapcore.NewCore(consoleEncoder, syncerStdout, lvlWarn),
			zapcore.NewCore(consoleEncoder, syncerStderr, lvlError),
		)
	}

	// NewTee
	core := zapcore.NewTee(cores...)
	log := zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(skip))

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := []zap.Field{}
		for k, v := range l.opts.Fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	// Adding namespace
	if namespace, ok := l.opts.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	//l.cfg = zapConfig
	l.zap = log
	l.fields = make(map[string]interface{})

	return nil
}

func (l *zaplog) Fields(fields map[string]interface{}) logger.Logger {
	l.Lock()
	nfields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		nfields[k] = v
	}
	l.Unlock()
	for k, v := range fields {
		nfields[k] = v
	}

	data := make([]zap.Field, 0, len(nfields))
	for k, v := range fields {
		data = append(data, zap.Any(k, v))
	}

	zl := &zaplog{
		cfg:    l.cfg,
		zap:    l.zap.With(data...),
		opts:   l.opts,
		fields: make(map[string]interface{}),
	}

	return zl
}

func (l *zaplog) Error(err error) logger.Logger {
	return l.Fields(map[string]interface{}{"error": err})
}

func (l *zaplog) Log(level logger.Level, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprint(args...)
	switch lvl {
	case TraceLevel:
		// Trace
		l.trace(msg, data...)
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) Logf(level logger.Level, format string, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprintf(format, args...)
	switch lvl {
	case TraceLevel:
		// Trace
		l.trace(msg, data...)
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) trace(msg string, fields ...zapcore.Field) {
	// Output directly
	ce := l.zap.Check(TraceLevel, msg)
	if ce != nil {
		ce.Write(fields...)
	}
}

func (l *zaplog) String() string {
	return "zlogger"
}

func (l *zaplog) Options() logger.Options {
	return l.opts
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...logger.Option) (logger.Logger, error) {
	// Default options
	options := logger.Options{
		Level:   logger.InfoLevel,
		Fields:  make(map[string]interface{}),
		Out:     os.Stderr,
		Context: context.Background(),
	}

	l := &zaplog{opts: options}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}

	return l, nil
}

func loggerToZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.TraceLevel:
		return TraceLevel
	case logger.DebugLevel:
		return zap.DebugLevel
	case logger.InfoLevel:
		return zap.InfoLevel
	case logger.WarnLevel:
		return zap.WarnLevel
	case logger.ErrorLevel:
		return zap.ErrorLevel
	case logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func zapToLoggerLevel(level zapcore.Level) logger.Level {
	switch level {
	case TraceLevel:
		return logger.TraceLevel
	case zap.DebugLevel:
		return logger.DebugLevel
	case zap.InfoLevel:
		return logger.InfoLevel
	case zap.WarnLevel:
		return logger.WarnLevel
	case zap.ErrorLevel:
		return logger.ErrorLevel
	case zap.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

// Options definations
type Options struct {
	logger.Options
}

type callerSkipKey struct{}

// WithCallerSkip : CallerSkip
func WithCallerSkip(i int) logger.Option {
	return logger.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) logger.Option {
	return logger.SetOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) logger.Option {
	return logger.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

// WithNamespace sets namespace for logger
func WithNamespace(namespace string) logger.Option {
	return logger.SetOption(namespaceKey{}, namespace)
}

type basePathKey struct{}

// WithBasePath sets rolling base path
func WithBasePath(basePath string) logger.Option {
	return logger.SetOption(basePathKey{}, basePath)
}

type echoKey struct{}

// WithEcho sets echo back
func WithEcho(echo bool) logger.Option {
	return logger.SetOption(echoKey{}, echo)
}

type diskKey struct{}

// WithDisk sets write disk
func WithDisk(disk bool) logger.Option {
	return logger.SetOption(diskKey{}, disk)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
