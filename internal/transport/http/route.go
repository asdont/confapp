package http

import (
	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	swagGin "github.com/swaggo/gin-swagger"

	"confapp/internal/app"
	"confapp/internal/handler"
)

func setRoutes(tools *app.Tools, router *gin.Engine) {
	router.GET("/doc/*any", swagGin.WrapHandler(swagFiles.Handler))

	v1GroupAPI := router.Group("/api")

	v1Group := v1GroupAPI.Group("/v1")
	{
		v1Group.POST("/config", handler.V1AddConfig(tools))
		v1Group.POST("/config/update", handler.V1AddNewVersionConfig(tools))

		v1Group.PUT("/config", handler.V1UpdateConfig(tools))
		v1Group.PUT("/config/any", handler.V1UpdateAllConfigs(tools))

		v1Group.GET("/config", handler.V1GetConfig(tools))
	}

	groupAuth := router.Group("/api", gin.BasicAuth(gin.Accounts{
		tools.Conf.Server.DeletionOperationsLogin: tools.Conf.Server.DeletionOperationsPass,
	}))

	v1GroupAuth := groupAuth.Group("/v1")
	{
		v1GroupAuth.DELETE("/service", handler.V1DeleteService(tools))
		v1GroupAuth.DELETE("/config", handler.V1DeleteConfig(tools))
	}
}
