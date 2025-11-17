package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Infof(format string, args any)
	Errorf(format string, args any)
}

type zapLogger struct {
	appLogger    *zap.Logger
	accessLogger *zap.Logger
}

func NewLogger() *zapLogger {
	// --- app error logger ---
	errorWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	})

	// --- access logger ---
	accessWS := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/access.log",
		MaxSize:    20,
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
	})

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// App logger: hanya untuk ERROR/INFO internal
	appCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		errorWS,
		zap.InfoLevel,
	)

	// Access logger: untuk request
	accessCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		accessWS,
		zap.InfoLevel,
	)

	return &zapLogger{
		appLogger:    zap.New(appCore, zap.AddCaller()),
		accessLogger: zap.New(accessCore),
	}
}

func (l *zapLogger) Infof(format string, args any) {
	l.appLogger.Sugar().Infof(format, args)
}

func (l *zapLogger) Errorf(format string, args any) {
	l.appLogger.Sugar().Errorf(format, args)
}

// ========== G I N  M I D D L E W A R E  ==========
func GinAccessLogger(l *zapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()

		l.accessLogger.Info("access log",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
