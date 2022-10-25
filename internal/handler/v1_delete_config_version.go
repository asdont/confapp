package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

type DeleteConfigQuery struct {
	Service string `form:"service" binding:"required" example:"managed-k8s"`
	Version int    `form:"v" binding:"required" example:"24"`
}

// V1DeleteConfig удаляет конфиг по имени сервиса и версии конфига(помечает удалённым, удаляет через 90 дней).
//
// @Summary удалить версию конфига по названию сервиса(помечает удалённым, удаляет через 90 дней)
// @Tags delete
// @Produce json
// @Param service query string true "Название сервиса удаляемого конфига"
// @Param v query integer true "Номер версии"
// @Success 200 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config [delete]
func V1DeleteConfig(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q DeleteConfigQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		tools.Mu.Lock()
		defer tools.Mu.Unlock()

		if err := model.DeleteConfigVersion(tools, q.Service, q.Version); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusOK, HTTPStatus{
			Status:  StatusDeleted,
			Version: q.Version,
		})
	}
}
