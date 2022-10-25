package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

type DeleteServiceQuery struct {
	Service string `form:"service" binding:"required" example:"managed-k8s"`
}

// V1DeleteService удаляет сервис со всеми версиями конфигов(помечает удалённым все версии, удаляет через 90 дней).
//
// @Summary удалить сервис по названию(помечает удалёнными все версии, удаляет через 90 дней)
// @Tags delete
// @Produce json
// @Param service query string true "Имя удаляемого сервиса"
// @Success 200 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/service [delete]
func V1DeleteService(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q DeleteServiceQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		tools.Mu.Lock()
		defer tools.Mu.Unlock()

		if err := model.DeleteServiceEasy(tools, q.Service); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusOK, HTTPStatus{
			Status:  StatusDeleted,
			Version: -1,
		})
	}
}
