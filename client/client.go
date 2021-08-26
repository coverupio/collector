package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient  *http.Client
	environment string
}

type Option func(c *Client)

func WithHTTPClient(cl *http.Client) Option {
	return func(c *Client) {
		c.httpClient = cl
	}
}

func WithEnvironment(env string) Option {
	return func(c *Client) {
		c.environment = env
	}
}

func New(opts ...Option) (*Client, error) {
	cl := &Client{
		environment: "production",
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
	for _, opt := range opts {
		opt(cl)
	}
	return cl, nil
}

func (c *Client) UploadCoverage(apiKey, repo, commit, format string, data io.Reader) error {
	covData, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read coverage data: %w", err)
	}

	url := apiURLForEnv(c.environment, "coverage.WriteCoverage")
	request := struct {
		APIKey       string `json:"apiKey"`
		Commit       string `json:"commit"`
		Repository   string `json:"repository"`
		ReportFormat string `json:"reportFormat"`
		Coverage     []byte `json:"coverage"`
	}{
		APIKey:       apiKey,
		Commit:       commit,
		Repository:   repo,
		ReportFormat: format,
		Coverage:     covData,
	}
	reqData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		return fmt.Errorf("failed to send coverage data to coverup.io: %w", err)
	}
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to send coverage data to coverup.io. Response %d: %s", res.StatusCode, string(body))
	}
	return nil
}

func apiURLForEnv(env, method string) string {
	switch env {
	case "local":
		return fmt.Sprintf("http://localhost:4060/%s", method)
	default:
		return fmt.Sprintf("https://coverup-enc-dr82.encoreapi.com/%s/%s", env, method)
	}
}
