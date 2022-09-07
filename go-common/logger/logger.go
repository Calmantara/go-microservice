//go:generate mockgen -source logger.go -destination mock/logger_mock.go -package mock

package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConst string

func (l LoggerConst) String() string {
	return string(l)
}

const (
	CorrelationKey LoggerConst = "Correlation-ID"
)

type Option func(z *zap.Config)

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

type CustomLogger interface {
	WithContext(ctx context.Context) *zap.SugaredLogger
	Logger() *zap.SugaredLogger
}

type CustomLoggerImpl struct {
	sugar *zap.SugaredLogger
}

func NewWrappedZapLogger(ops ...Option) *zap.SugaredLogger {
	// logger setup
	cfg := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		Development:       false,
		DisableStacktrace: true,
	}
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.EncodeLevel = customLevelEncoder

	//iterate all option function
	for _, v := range ops {
		v(&cfg)
	}

	logger, _ := cfg.Build()
	logger.Info("Initialization zap logger. . .")
	defer logger.Sync()

	return logger.Sugar()
}

func NewCustomLogger(ops ...Option) CustomLogger {
	sugar := NewWrappedZapLogger(ops...)
	defer sugar.Sync()
	return &CustomLoggerImpl{
		sugar: sugar,
	}
}

func (c *CustomLoggerImpl) WithContext(ctx context.Context) *zap.SugaredLogger {
	// check correlation ID is exist or not
	corr := ctx.Value(CorrelationKey.String())
	logger := c.sugar
	if corr != nil {
		logger = logger.With(CorrelationKey.String(), corr)
	}

	return logger
}

func (c *CustomLoggerImpl) Logger() *zap.SugaredLogger {
	// this function only return zap
	return c.sugar
}

func WithTimeKey(tk string) Option {
	return func(z *zap.Config) { z.EncoderConfig.TimeKey = tk }
}
func WithTimeFormat(encodeTime zapcore.TimeEncoder) Option {
	return func(z *zap.Config) { z.EncoderConfig.EncodeTime = encodeTime }
}
