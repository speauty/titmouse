package log

import (
	"github.com/golang-module/carbon"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	apiLog  *Log
	onceLog sync.Once
)

func Api() *Log {
	onceLog.Do(func() {
		apiLog = new(Log)
	})
	return apiLog
}

type Log struct { // 接入部分API
	logger  *zap.SugaredLogger
	cfg     *Cfg
	encoder zapcore.Encoder
}

func (customTL *Log) Init(cfg *Cfg) {
	if customTL.logger != nil {
		return
	}
	apiLog.cfg = cfg
	apiLog.initCfg()
	apiLog.genEncoder()
	apiLog.genLogger()
}

func (customTL *Log) Scope(scopeName string) *zap.SugaredLogger {
	if scopeName == "" {
		return customTL.logger
	}
	return customTL.logger.With(zap.String("scope", scopeName))
}

func (customTL *Log) Debug(args ...interface{}) {
	customTL.logger.Debug(args...)
}
func (customTL *Log) Info(args ...interface{}) {
	customTL.logger.Info(args...)
}
func (customTL *Log) Warn(args ...interface{}) {
	customTL.logger.Warn(args...)
}
func (customTL *Log) Error(args ...interface{}) {
	customTL.logger.Error(args...)
}
func (customTL *Log) Panic(args ...interface{}) {
	customTL.logger.Panic(args...)
}
func (customTL *Log) Fatal(args ...interface{}) {
	customTL.logger.Fatal(args...)
}

func (customTL *Log) initCfg() {
	defaultCfg := new(Cfg).Default()
	if apiLog.cfg == nil {
		apiLog.cfg = new(Cfg).Default()
	} else {
		if apiLog.cfg.InfoFile == "" {
			apiLog.cfg.InfoFile = defaultCfg.InfoFile
		}
		if apiLog.cfg.ErrorFile == "" {
			apiLog.cfg.ErrorFile = defaultCfg.ErrorFile
		}
	}
}

func (customTL *Log) genLogger() {
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.ErrorLevel
	})

	ioInfo := customTL.genWriter(customTL.cfg.InfoFile)
	ioErr := customTL.genWriter(customTL.cfg.ErrorFile)

	core := zapcore.NewTee(
		zapcore.NewCore(customTL.encoder, zapcore.AddSync(ioInfo), infoLevel),
		zapcore.NewCore(customTL.encoder, zapcore.AddSync(ioErr), errorLevel),
	)
	customTL.logger = zap.New(core, zap.AddCaller()).Sugar()
}

func (customTL *Log) genWriter(filename string) io.Writer {
	fd, err := rotatelogs.New(strings.Replace(filename, ".log", "", -1) + ".%Y%m%d.log")
	if err != nil {
		panic(err)
	}
	return fd
}

func (customTL *Log) genEncoder() {
	customTL.encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg", LevelKey: "level", TimeKey: "timestamp", StacktraceKey: "trace",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(carbon.FromStdTime(time).Layout(carbon.ShortDateTimeLayout))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}
