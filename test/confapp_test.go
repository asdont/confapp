package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-playground/assert/v2"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"confapp/internal/config"
	"confapp/internal/handler"
	"confapp/internal/postgres"
)

const (
	uriAddConfig        = "http://127.0.0.1:22952/api/v1/config"
	uriAddNewVersion    = "http://127.0.0.1:22952/api/v1/config/update"
	uriUpdateConfig     = "http://127.0.0.1:22952/api/v1/config"
	uriUpdateAllConfigs = "http://127.0.0.1:22952/api/v1/config/any"
	uriGetConfig        = "http://127.0.0.1:22952/api/v1/config"
	uriDeleteConfig     = "http://127.0.0.1:22952/api/v1/config"
)

const (
	testServiceName = "TEST_1"
)

func TestConfAppPositiveScenario(t *testing.T) {
	conf, err := config.GetConf("../config/conf.toml")
	require.NoError(t, err)

	db, err := postgres.CreatePool(conf)
	require.NoError(t, err)

	defer func(db *sql.DB) {
		err := db.Close()
		require.NoError(t, err)
	}(db)

	err = deleteTestDataFromDB(db, []string{testServiceName})
	require.NoError(t, err)

	client := &http.Client{
		Timeout: time.Second,
	}

	err = addConfig(t, client)
	require.NoError(t, err)

	err = addNewVersion(t, client)
	require.NoError(t, err)

	err = updateConfig(t, client)
	require.NoError(t, err)

	err = updateAllConfigs(t, client)
	require.NoError(t, err)

	err = getConfig(t, client)
	require.NoError(t, err)

	err = deleteConfig(t, client, conf.Server.DeletionOperationsLogin, conf.Server.DeletionOperationsPass)
	require.NoError(t, err)

	err = getDeletedConfig(t, client)
	require.NoError(t, err)

	err = deleteTestDataFromDB(db, []string{testServiceName})
	require.NoError(t, err)
}

func addConfig(t *testing.T, client *http.Client) error {
	data, err := json.Marshal(handler.AddConfigBody{
		Service: testServiceName,
		Data: []json.RawMessage{
			[]byte(`{"key1": "value1"}`),
			[]byte(`{"key2": "value2"}`),
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, uriAddConfig, bytes.NewReader(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handler.HTTPStatus
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, handler.HTTPStatus{
		Status:  handler.StatusSuccess,
		Version: 1,
	}, res)

	return nil
}

func addNewVersion(t *testing.T, client *http.Client) error {
	data, err := json.Marshal(handler.UpdateConfigBody{
		Service: testServiceName,
		Data: []json.RawMessage{
			[]byte(`{"key3": "value3"}`),
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, uriAddNewVersion, bytes.NewReader(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handler.HTTPStatus
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, handler.HTTPStatus{
		Status:  handler.StatusSuccess,
		Version: 2,
	}, res)

	return nil
}

func updateConfig(t *testing.T, client *http.Client) error {
	data, err := json.Marshal(handler.UpdateConfigBody{
		Service: testServiceName,
		Version: 1,
		Data: []json.RawMessage{
			[]byte(`{"key1": "updated_value"}`),
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, uriUpdateConfig, bytes.NewReader(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handler.HTTPStatus
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, handler.HTTPStatus{
		Status:  handler.StatusSuccess,
		Version: 1,
	}, res)

	return nil
}

func updateAllConfigs(t *testing.T, client *http.Client) error {
	data, err := json.Marshal(handler.UpdateAllConfigsBody{
		Service: testServiceName,
		Data: []json.RawMessage{
			[]byte(`{"key4": "value4"}`),
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, uriUpdateAllConfigs, bytes.NewReader(data))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handler.HTTPStatus
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, handler.HTTPStatus{
		Status: handler.StatusSuccess,
	}, res)

	return nil
}

func getConfig(t *testing.T, client *http.Client) error {
	uri, err := url.Parse(uriGetConfig)
	require.NoError(t, err)

	urlValues := &url.Values{
		"service": {testServiceName},
		"v":       {"1"},
	}

	uri.RawQuery = urlValues.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var params map[string]string
	err = json.Unmarshal(body, &params)
	require.NoError(t, err)

	requiredParams := map[string]string{
		"key1": "updated_value",
		"key2": "value2",
		"key4": "value4",
	}

	assert.Equal(t, requiredParams, params)

	return nil
}

func deleteConfig(t *testing.T, client *http.Client, login, pass string) error {
	uri, err := url.Parse(uriDeleteConfig)
	require.NoError(t, err)

	urlValues := &url.Values{
		"service": {testServiceName},
		"v":       {"1"},
	}

	uri.RawQuery = urlValues.Encode()

	req, err := http.NewRequest(http.MethodDelete, uri.String(), nil)
	require.NoError(t, err)

	req.SetBasicAuth(login, pass)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var res handler.HTTPStatus
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, handler.HTTPStatus{
		Status:  handler.StatusDeleted,
		Version: 1,
	}, res)

	return nil
}

func getDeletedConfig(t *testing.T, client *http.Client) error {
	uri, err := url.Parse(uriGetConfig)
	require.NoError(t, err)

	urlValues := &url.Values{
		"service": {testServiceName},
		"v":       {"1"},
	}

	uri.RawQuery = urlValues.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var params map[string]string
	err = json.Unmarshal(body, &params)
	require.NoError(t, err)

	requiredParams := map[string]string{}

	assert.Equal(t, requiredParams, params)

	return nil
}

func deleteTestDataFromDB(db *sql.DB, serviceNames []string) error {
	if _, err := db.Exec(`
		delete from
		    service
		where
		    name = any ($1)
	`,
		pq.Array(serviceNames),
	); err != nil {
		return err
	}

	return nil
}
