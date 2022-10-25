package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/natefinch/lumberjack.v2"

	"confapp/internal/config"
)

const logTimeFormat = "06-01-02 15:04:05"

func InitApp(fileName string, conf *config.Conf) *logrus.Logger {
	file := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    conf.Logger.AppMaxLogSizeMb,
		MaxBackups: conf.Logger.AppMaxNumOfBackups,
		MaxAge:     conf.Logger.AppMaxBackupAgeDay,
		Compress:   true,
	}

	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, file),
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: logTimeFormat,
			LogFormat:       "%time% [%lvl%] %msg%\n",
		},
	}

	return logger
}

func InitServer(fileName string, conf *config.Conf) *lumberjack.Logger {
	logger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    conf.Logger.ServerMaxLogSizeMb,
		MaxBackups: conf.Logger.ServerMaxNumOfBackups,
		MaxAge:     conf.Logger.ServerAppMaxBackupAgeDay,
		Compress:   true,
	}

	return logger
}
