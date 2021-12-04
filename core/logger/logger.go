package logger

import (
	"fiber/core/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化Logger
func Init() {
	writeSyncer := getLogWriter(config.Viper.GetString("logger.path")+"/log/log.log", config.Viper.GetInt("logger.maxSize"))
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  filename,
		MaxSize:   maxSize,
		LocalTime: true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
