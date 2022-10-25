package app

import (
	"database/sql"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"confapp/internal/config"
)

type Tools struct {
	Mu     *sync.RWMutex
	Conf   *config.Conf
	DB     *sql.DB
	Logger *Loggers
}

type Loggers struct {
	App    *logrus.Logger
	Server *lumberjack.Logger
}
