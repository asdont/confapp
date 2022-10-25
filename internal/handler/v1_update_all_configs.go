package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

type UpdateAllConfigsBody struct {
	Service string            `json:"service" binding:"required"`
	Data    []json.RawMessage `json:"data" binding:"required"`
}

// V1UpdateAllConfigs обновляет и/или добавляет новые поля во всех версиях конфигов сервиса.
//
// @Summary обновить и/или добавить новые поля во все версии конфига сервиса
// @Tags update
// @Accept json
// @Produce json
// @Param data body string false "Заменить кавычки на двойные - {'service': 'managed-k8s', 'data': [{'key1': 'changed_value'}]}"
// @Success 200 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config/any [put]
func V1UpdateAllConfigs(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b UpdateAllConfigsBody
		if err := c.ShouldBind(&b); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		params, err := extractDataFromJSON(b.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		// Вместо изоляции транзакции.
		tools.Mu.Lock()
		defer tools.Mu.Unlock()

		serviceConfigs, err := model.GetServiceConfigs(tools, b.Service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		if len(serviceConfigs) == 0 {
			c.JSON(http.StatusBadRequest, HTTPStatus{
				Error: StatusVersionsNotFound},
			)

			return
		}

		mergeConfigs(params, serviceConfigs)

		if err := model.UpdateAllConfigs(tools, serviceConfigs); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusOK, HTTPStatus{
			Status: StatusSuccess,
		})
	}
}

func mergeConfigs(params map[string]string, serviceConfigs map[int]map[string]string) {
	for _, versionConfigs := range serviceConfigs {
		for param, value := range params {
			versionConfigs[param] = value
		}
	}
}
