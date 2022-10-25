package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"confapp/internal/app"
	"confapp/internal/model"
)

var (
	errVersionsEqual       = errors.New("version are equal")
	errLastVersionNotFound = errors.New("last version not found")
)

type AddNewVersionConfigBody struct {
	Service string            `json:"service" binding:"required" example:"managed-k8s"`
	Data    []json.RawMessage `json:"data" binding:"required"`
}

// V1AddNewVersionConfig обновляет последнюю версию конфига.
//
// @Summary обновить последний конфиг
// @Description Создаёт копию последней версии с новыми и/или измененными параметрами, под новым номером версии
// @Tags create
// @Accept json
// @Produce json
// @Param data body string false "Заменить кавычки на двойные - {'service': 'managed-k8s', 'data': [{'key3': 'value3'}]}"
// @Success 201 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config/update [post]
func V1AddNewVersionConfig(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b AddNewVersionConfigBody
		if err := c.ShouldBind(&b); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		params, err := extractDataFromJSON(b.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		// Вместо изоляции транзакции
		tools.Mu.Lock()
		defer tools.Mu.Unlock()

		lastVersionNumber, paramsLastVersion, err := model.GetLastVersionConfig(tools, b.Service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		if lastVersionNumber == 0 {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: errLastVersionNotFound.Error()})

			return
		}

		if err := mergeParameters(params, paramsLastVersion); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: fmt.Sprintf("latest and this: %v", err)})

			return
		}

		if err := model.AddNewVersionConfig(tools, b.Service, lastVersionNumber, paramsLastVersion); err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusCreated, HTTPStatus{
			Status:  StatusSuccess,
			Version: lastVersionNumber + 1,
		})
	}
}

func mergeParameters(params, paramsLatestVersion map[string]string) error {
	duplicateCounter := 0
	for param, value := range params {
		if val, exist := paramsLatestVersion[param]; exist && value == val {
			duplicateCounter++
		}

		paramsLatestVersion[param] = value
	}

	if duplicateCounter == len(params) {
		return errVersionsEqual
	}

	return nil
}
