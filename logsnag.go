package logsnag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogResponse struct {
	Data interface{} `json:"data"`
}

type PublishPayload struct {
	Project     string                 `json:"project"`
	Channel     string                 `json:"channel"`
	Event       string                 `json:"event"`
	Description string                 `json:"description,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
	Notify      bool                   `json:"notify,omitempty"`
	Parser      map[string]interface{} `json:"parser,omitempty"`
}

type IPublishPayloadOptions struct {
	Description string                 `json:"description,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
	Notify      bool                   `json:"notify,omitempty"`
	Parser      map[string]interface{} `json:"parser,omitempty"`
}

type InsightPayload struct {
	Project string      `json:"project"`
	Title   string      `json:"title"`
	Value   interface{} `json:"value"`
	Icon    string      `json:"icon,omitempty"`
}
type InsightPayloadOptions struct {
	Icon string `json:"icon,omitempty"`
}

type logsnagConfig struct {
	ProjectName string
	Client      APIClient
}

type Logsnag interface {
	Publish(channel string, event string, options IPublishPayloadOptions) (LogResponse, error)
	Insight(title string, value interface{}, options InsightPayloadOptions) (resp LogResponse, err error)
}

func NewLogsnag(project string, cfg *APIClient) Logsnag {
	return &logsnagConfig{
		ProjectName: project,
		Client:      *NewApiClient(cfg.Token, &Config{}),
	}
}

func (l *logsnagConfig) Publish(channel string, event string, options IPublishPayloadOptions) (resp LogResponse, err error) {
	endpoint := fmt.Sprintf(l.Client.Config.BaseBackendUrl+"/%s", "log")

	// Make post request to endpoint
	payload := PublishPayload{
		Project:     l.ProjectName,
		Channel:     channel,
		Event:       event,
		Description: options.Description,
		Icon:        options.Icon,
		Notify:      options.Notify,
		Tags:        options.Tags,
		Parser:      options.Parser,
	}

	fmt.Println(payload)
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))

	if err != nil {
		return resp, err
	}

	resp.Data, err = l.Client.SendRequest(req)
	return resp, err
}
func (l *logsnagConfig) Insight(title string, value interface{}, options InsightPayloadOptions) (resp LogResponse, err error) {
	endpoint := fmt.Sprintf(l.Client.Config.BaseBackendUrl+"/%s", "insight")

	// Make post request to endpoint
	jsonBody, _ := json.Marshal(InsightPayload{
		Project: l.ProjectName,
		Title:   title,
		Value:   value,
		Icon:    options.Icon,
	})

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))

	if err != nil {
		return resp, err
	}
	resp.Data, err = l.Client.SendRequest(req)
	return resp, err
}
