package lib

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/RealLiuSha/echo-admin/constants"
	"github.com/RealLiuSha/echo-admin/pkg/file"
)

// Zap SugaredLogger by default
// DesugarZap performance-sensitive code
type Logger struct {
	Zap        *zap.SugaredLogger
	DesugarZap *zap.Logger
}

func NewLogger(config Config) Logger {
	var options []zap.Option
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     localTimeEncoder,
	}

	if config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level := zap.NewAtomicLevelAt(toLevel(config.Log.Level))

	core := zapcore.NewCore(encoder, toWriter(config), level)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.WarnLevel)
	options = append(options,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(core, options...)
	return Logger{Zap: logger.Sugar(), DesugarZap: logger}
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(constants.TimeFormat))
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func toWriter(config Config) zapcore.WriteSyncer {
	fp := ""
	sp := string(filepath.Separator)

	fp, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
	fp += sp + "logs" + sp

	if config.Log.Directory != "" {
		if err := file.EnsureDirRW(config.Log.Directory); err != nil {
			fp = config.Log.Directory
		}
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(fp, config.Name) + ".log",
			MaxSize:    100,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   true,
		}),
	)
}
