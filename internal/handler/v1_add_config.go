package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/perimeterx/marshmallow"

	"confapp/internal/app"
	"confapp/internal/model"
)

var (
	errDuplicate       = errors.New("duplicate")
	errIncorrectFormat = errors.New("incorrect format")
)

type AddConfigBody struct {
	Service string            `json:"service" binding:"required" example:"managed-k8s"`
	Data    []json.RawMessage `json:"data" binding:"required"`
}

// V1AddConfig создаёт конфиг для нового сервиса.
//
// @Summary создать конфиг для нового сервиса
// @Tags create
// @Accept json
// @Produce json
// @Param data body string false "Заменить кавычки на двойные - {'service': 'managed-k8s', 'data': [{'key1': 'value1'}, {'key2': 'value2'}]}"
// @Success 201 {object} HTTPStatus
// @Failure 400 {object} HTTPStatus
// @Failure 500 {object} HTTPStatus
// @Router /v1/config [post]
func V1AddConfig(tools *app.Tools) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b AddConfigBody
		if err := c.ShouldBind(&b); err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		params, err := extractDataFromJSON(b.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, HTTPStatus{Error: err.Error()})

			return
		}

		versionNumber, err := model.AddConfig(tools, b.Service, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, HTTPStatus{Error: err.Error()})

			return
		}

		c.JSON(http.StatusCreated, HTTPStatus{
			Status:  StatusSuccess,
			Version: versionNumber,
		})
	}
}

// Извлекает ключи и значения из списка структур, типа:
// [
//
//	{"key1": "value1"},
//	{"key2": "value2"}
//
// ].
func extractDataFromJSON(data []json.RawMessage) (map[string]string, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty: %w", errIncorrectFormat)
	}

	params := make(map[string]string)

	for _, d := range data {
		kv, err := marshmallow.Unmarshal(d, &struct{}{})
		if err != nil {
			return nil, fmt.Errorf("unmarshal: %w: %v", err, errIncorrectFormat)
		}

		for k, v := range kv {
			vString, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("unmarshal string %v: %w", v, errIncorrectFormat)
			}

			if _, exist := params[k]; exist {
				return nil, fmt.Errorf("%s: %w", k, errDuplicate)
			}

			if strings.Contains(k, " ") {
				return nil, fmt.Errorf("%w: space is not allowed: <%s>", errIncorrectFormat, k)
			}

			params[k] = vString
		}
	}

	return params, nil
}
