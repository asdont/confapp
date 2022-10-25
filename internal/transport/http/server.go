package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
)

const (
	maxHeaderBytes              = 1 << 20
	maxServerShutdownTimeSecond = 30
)

func RunServerHTTP(tools *app.Tools, chStop chan struct{}) error {
	gin.DisableConsoleColor()

	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.MultiWriter(tools.Logger.Server, os.Stdout)

	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.LoggerWithWriter(gin.DefaultWriter),
	)

	setRoutes(tools, router)

	server := &http.Server{
		Addr:           ":" + tools.Conf.Server.Port,
		Handler:        router,
		ReadTimeout:    time.Second * time.Duration(tools.Conf.Server.ReadTimeoutSec),
		WriteTimeout:   time.Second * time.Duration(tools.Conf.Server.WriteTimeoutSec),
		MaxHeaderBytes: maxHeaderBytes,
	}

	tools.Logger.App.Infof("...server port: %s", tools.Conf.Server.Port)

	chErr := make(chan error)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			chErr <- fmt.Errorf("http server: %w", err)
		}
	}()

	select {
	case err := <-chErr:
		return err

	case <-chStop:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*maxServerShutdownTimeSecond)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("[SERVER][FAIL] forced stop: %w", err)
		}

		tools.Logger.App.Infof("[SERVER][OK] stop")
	}

	return nil
}
