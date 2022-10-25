package main

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"

	_ "confapp/docs"

	"confapp/internal/app"
	"confapp/internal/config"
	"confapp/internal/logger"
	"confapp/internal/model"
	"confapp/internal/postgres"
	"confapp/internal/transport/http"
)

const (
	fileConf      = "config/conf.toml"
	fileLogApp    = "logs/app.log"
	fileLogServer = "logs/server.log"
)

// @title API ConfApp
// @version 1.0

// @contact.name Company
// @contact.url https://example.ru/
// @contact.email example@mail.ru

// @host 127.0.0.1:22952
// @schemes http
// @BasePath /api
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth
func main() {
	var mu sync.RWMutex

	conf, err := config.GetConf(fileConf)
	if err != nil {
		log.Fatal(err)
	}

	postgresPool, err := postgres.CreatePool(conf)
	if err != nil {
		log.Fatal(err)
	}

	loggerApp := logger.InitApp(fileLogApp, conf)
	loggerServer := logger.InitServer(fileLogServer, conf)

	tools := &app.Tools{
		Mu:   &mu,
		Conf: conf,
		DB:   postgresPool,
		Logger: &app.Loggers{
			App:    loggerApp,
			Server: loggerServer,
		},
	}

	chStop := make(chan struct{})

	go func() {
		if err := appMaintenance(tools, chStop); err != nil {
			tools.Logger.App.Fatal(err)
		}
	}()

	if err := http.RunServerHTTP(tools, chStop); err != nil {
		tools.Logger.App.Fatalf("server http: %v", err)
	}
}

func appMaintenance(tools *app.Tools, chStop chan struct{}) error {
	chErr := make(chan error)

	go app.SysStop(chStop)

	go func() {
		if err := model.DeleteOldRowsHard(tools); err != nil {
			chErr <- fmt.Errorf("delete old rows hard: %w", err)
		}
	}()

	return <-chErr
}
