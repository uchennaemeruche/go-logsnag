package logsnag

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewApiClient(t *testing.T) {
	logsnagBearerToken := "YOUR_TOKEN"
	baseBackendURL := "https://api.logsnag.com/v1"

	// Test NewApiClient function
	t.Run("NewApiClient function should return an APIClient object", func(t *testing.T) {
		client := NewApiClient(logsnagBearerToken, &Config{})
		assert.IsType(t, &APIClient{}, client)
		assert.Equal(t, logsnagBearerToken, client.Token)
		assert.Equal(t, baseBackendURL, client.Config.BaseBackendUrl)
		assert.IsType(t, &http.Client{}, client.Config.HttpClient)
	})

	// Test NewApiClient function with custom BaseBackendUrl
	customBaseBackendURL := "https://custom.api.logsnag.com/v1"
	t.Run("NewApiClient function should use custom BaseBackendUrl if provided", func(t *testing.T) {
		client := NewApiClient(logsnagBearerToken, &Config{BaseBackendUrl: customBaseBackendURL})
		assert.Equal(t, customBaseBackendURL, client.Config.BaseBackendUrl)
	})

	// Test NewApiClient function with custom HttpClient
	customHttpClient := &http.Client{Timeout: 30 * time.Second}
	t.Run("NewApiClient function should use custom HttpClient if provided", func(t *testing.T) {
		client := NewApiClient(logsnagBearerToken, &Config{HttpClient: customHttpClient})
		assert.Equal(t, customHttpClient, client.Config.HttpClient)
	})
}

func TestSendRequest(t *testing.T) {
	logsnagBearerToken := "YOUR_TOKEN"

	t.Run("SendRequest function should send the request and return a response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"message": "Event published successfully"}`))
		}))
		defer server.Close()

		client := NewApiClient(logsnagBearerToken, &Config{BaseBackendUrl: server.URL})
		jsonBody, _ := json.Marshal(PublishPayload{
			Project:     "project_name",
			Channel:     "test_channel",
			Event:       "test_event",
			Description: "test_description",
			Notify:      true,
		})
		req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/log", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)

		v, err := client.SendRequest(req)
		assert.NoError(t, err)

		data, ok := v.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Event published successfully", data["message"])
	})

	t.Run("SendRequest function should return an error if request fails", func(t *testing.T) {
		client := NewApiClient(logsnagBearerToken, &Config{BaseBackendUrl: "http://localhost:1234"})
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:1234/log", nil)
		_, err := client.SendRequest(req)
		assert.Error(t, err)
	})
}

func TestPublish(t *testing.T) {
	logsnagBearerToken := "YOUR_TOKEN"
	logsnagProject := "project_name"

	options := &PublishPayload{
		Project:     "project_name",
		Channel:     "channel",
		Event:       "event",
		Description: "description",
		Notify:      true,
	}

	t.Run("Publish method should send a request and return a response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var reqBody PublishPayload
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				log.Printf("error in unmarshalling %+v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			t.Run("Header must contain Bearer token", func(t *testing.T) {
				bearerToken := r.Header.Get("Authorization")
				assert.True(t, strings.HasPrefix(bearerToken, "Bearer"))
				assert.True(t, strings.Contains(bearerToken, logsnagBearerToken))
			})

			// Test the URL and request method
			t.Run("URL and Request method must not be null", func(t *testing.T) {
				expectedURL := "/v1/log"
				assert.Equal(t, expectedURL, r.RequestURI)
				assert.Equal(t, http.MethodPost, r.Method)
			})

			w.Write([]byte(`{
				"channel":"channel",
				"description":"description",
				"event":"event",
				"notify":true,
				"parser":"text",
				"project":"project_name"
			}`))

			t.Run("Request body must match expected", func(t *testing.T) {
				assert.Equal(t, options, &reqBody)
			})

		}))

		defer server.Close()

		client := NewApiClient(logsnagBearerToken, &Config{BaseBackendUrl: server.URL})

		// logsnag := NewLogsnag(&NewLogParams{"", logsnagProject, client})
		logsnag := NewLogsnag(logsnagProject, client)
		resp, err := logsnag.Publish(
			options.Channel, options.Event, IPublishPayloadOptions{Description: options.Description, Notify: options.Notify},
		)
		assert.NoError(t, err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, options.Channel, data["channel"])
		assert.Equal(t, options.Event, data["event"])
		assert.Equal(t, options.Project, data["project"])
	})

	t.Run("Publish method should return an error if request fails", func(t *testing.T) {
		// client := NewApiClient(logsnagBearerToken, &Config{})
		logsnag := NewLogsnag(logsnagProject, &APIClient{})
		_, err := logsnag.Publish("channel", "event", IPublishPayloadOptions{Description: "description", Notify: true})
		assert.Error(t, err)
	})
}
func TestInsight(t *testing.T) {
	logsnagBearerToken := "YOUR_TOKEN"
	logsnagProject := "project_name"

	options := &InsightPayload{
		Project: "project_name",
		Title:   "title",
		Value:   "new_insight_value",
	}

	t.Run("Test Insight", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Run("Header must contain Bearer token", func(t *testing.T) {
				bearerToken := r.Header.Get("Authorization")
				assert.True(t, strings.HasPrefix(bearerToken, "Bearer"))
				assert.True(t, strings.Contains(bearerToken, logsnagBearerToken))
			})

			// Test the URL and request method
			t.Run("URL and Request method must not be null", func(t *testing.T) {
				expectedURL := "/v1/insight"
				assert.Equal(t, expectedURL, r.RequestURI)
				assert.Equal(t, http.MethodPost, r.Method)
			})

			w.Write([]byte(`{
				"title":"title",
				"value":"new_insight_value"
			}`))
		}))

		defer server.Close()

		client := NewApiClient(logsnagBearerToken, &Config{BaseBackendUrl: server.URL})

		logsnag := NewLogsnag(logsnagProject, client)
		resp, err := logsnag.Insight(
			options.Title, options.Value, InsightPayloadOptions{},
		)
		assert.NoError(t, err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, options.Title, data["title"])
		assert.Equal(t, options.Value, data["value"])
	})
}
