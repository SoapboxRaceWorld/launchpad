package launchpad

import (
	"go.uber.org/zap"
	"fmt"
	"go.uber.org/zap/zapcore"
)

var (
	mnLogger *zap.SugaredLogger
)

func GetLogger() *zap.SugaredLogger {
	if mnLogger == nil {
		var level zapcore.Level

		if isDebug {
			level = zap.DebugLevel
		} else {
			level = zap.InfoLevel
		}

		config := zap.Config{
			Level:            zap.NewAtomicLevelAt(level),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, err := config.Build()

		if err != nil {
			fmt.Println(err)
			panic("Can't create logger")
		}

		mnLogger = logger.Sugar()
	}

	return mnLogger
}
