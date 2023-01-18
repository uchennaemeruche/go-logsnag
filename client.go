package logsnag

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	BaseBackendUrl string
	HttpClient     *http.Client
}

type APIClient struct {
	Token  string
	Config *Config
}

func NewApiClient(token string, config *Config) *APIClient {
	config.BaseBackendUrl = buildBaseUrl(config)

	if config.HttpClient == nil {
		config.HttpClient = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	client := &APIClient{
		Token: token,
	}
	client.Config = config

	return &APIClient{
		Token:  token,
		Config: config,
	}
}

func (client *APIClient) SendRequest(req *http.Request) (v interface{}, err error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))

	res, err := client.Config.HttpClient.Do(req)
	if err != nil {
		return v, errors.Wrap(err, "failed to execute request")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode >= http.StatusMultipleChoices {
		return v, errors.Errorf(
			`status code %d, %s`, res.StatusCode,
			string(body),
		)
	}

	if err = json.Unmarshal(body, &v); err != nil {
		return v, errors.Wrap(err, "unable to unmarshal response body")
	}
	return v, nil
}

func buildBaseUrl(config *Config) string {
	apiVersion := "v1"
	if config.BaseBackendUrl == "" {
		return fmt.Sprintf("https://api.logsnag.com/%s", apiVersion)
	}
	if strings.Contains(config.BaseBackendUrl, "logsnag.com/v") {
		return config.BaseBackendUrl
	}
	return fmt.Sprintf(config.BaseBackendUrl+"/%s", apiVersion)
}
