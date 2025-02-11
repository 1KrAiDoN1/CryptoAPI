package coincap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		return nil, errors.New("timeout can't be zero")

	}
	return &Client{
		client: &http.Client{
			Timeout: timeout,
			Transport: &LoggingRoundTripper{
				logger: os.Stdout,
				next:   http.DefaultTransport,
			},
		},
	}, nil

}

func (c Client) GetAsset(name string) (Asset, error) {
	url := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", name)
	resp, err := c.client.Get(url)
	if err != nil {
		return Asset{}, err
	}
	defer resp.Body.Close()

	// Чтение содержимого ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, err
	}
	var r AssetResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return Asset{}, err
	}
	return r.Asset, nil
}
