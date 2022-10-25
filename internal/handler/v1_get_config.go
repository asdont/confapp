package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

type GetConfigQuery struct {
	Service string `form:"service" binding:"required" example:"managed-k8s"`
	Version int    `form:"v" binding:"required" example:"24"`
}

// V1GetConfig получает конфиг по названию сервиса и номеру версии конфига.
//
// @Summary получить конфиг по названию сервиса и номеру версии конфига
// @Tags get
// @Produce json
// @Param service query string true "Название сервиса удаляемого конфига"
// @Param v query integer true "Номер версии"
// @Success 200 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config [get]
func V1GetConfig(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q GetConfigQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		params, err := model.GetConfig(tools, q.Service, q.Version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusOK, params)
	}
}
