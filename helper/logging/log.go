package logging

import (
	"campus/helper/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var log *zap.Logger

func Setup() {
	var err error
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	currLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg := zap.Config{
		Level:         currLevel,
		Development:   true,
		Encoding:      "json",
		EncoderConfig: encoderCfg,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPaths: []string{
			"stdout", // 同时输出到终端  不需要则删除
			setting.AppSetting.LogSavePath,
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	log,err = cfg.Build()
	if err != nil {
		println(err.Error())
		return
	}
	log = log.Named("log")

}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
func Close() {
	log.Sync()
}

func Info(msg string,fields ...zap.Field) {
	log.Info(msg,fields...)
}
func Warn(msg string, fields ... zap.Field) {
	log.Warn(msg,fields...)
}
func WarnMsg(msg string, err error) {
	log.Warn(msg,zap.String("error",err.Error()))
}
func Error(msg string,fields ...zap.Field) {
	log.Error(msg,fields...)
}
func ErrorMsg(msg string, err error) {
	log.Error(msg,zap.String("error",err.Error()))
}

func DatabaseError(err error)  {
	log.Error("数据库连接错误",zap.String("error",err.Error()))
}
