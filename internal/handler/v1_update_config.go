package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

type UpdateConfigBody struct {
	Service string            `json:"service" binding:"required"`
	Version int               `json:"v" binding:"required"`
	Data    []json.RawMessage `json:"data" binding:"required"`
}

// V1UpdateConfig обновляет конфиг по названию сервиса и номеру версии.
//
// @Summary обновить одну версию конфига по названию сервиса и номеру версии
// @Tags update
// @Accept json
// @Produce json
// @Param data body string false "Заменить кавычки на двойные - {'service': 'managed-k8s', 'v': 224, 'data': [{'key1': 'value1-1'}]}"
// @Success 200 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config [put]
func V1UpdateConfig(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b UpdateConfigBody
		if err := c.ShouldBind(&b); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		params, err := extractDataFromJSON(b.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		if err := model.UpdateConfig(tools, b.Service, b.Version, params); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, HTTPStatus{
					Error: StatusVersionNotFound,
				})

				return
			}

			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusOK, HTTPStatus{
			Status:  StatusSuccess,
			Version: b.Version,
		})
	}
}
